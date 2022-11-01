package surl

import (
	"crypto/subtle"
	"errors"
	"hash"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/blake2b"
)

var (
	// ErrInvalidSignature is returned when the provided token's
	// signature is not valid.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidSignedURL is returned when the signed URL's format is
	// invalid.
	ErrInvalidSignedURL = errors.New("invalid signed URL")
	// ErrExpired is returned by when the signed URL's expiry has been
	// exceeded.
	ErrExpired = errors.New("URL has expired")
	// Default formatter is the query formatter.
	DefaultFormatter = WithQueryFormatter()
)

// Signer is capable of signing and verifying signed URLs with an expiry.
type Signer struct {
	mu        sync.Mutex
	hash      hash.Hash
	dirty     bool
	skipQuery bool
	prefix    string

	Formatter
}

// New constructs a new signer, performing the one-off task of generating a
// secure hash from the key. The key must be between 0 and 64 bytes long;
// anything longer is stripped off.
func New(key []byte, opts ...Option) *Signer {
	hash, err := blake2b.New256(key)
	if err != nil {
		// The only possible error that can be returned here is if the key
		// is larger than 64 bytes - which the blake2b hash will not accept.
		// This is a case that is so easily avoidable when using this package
		// and since chaining is convenient for this package.  We're going
		// to do the below to handle this possible case so we don't have
		// to return an error.
		hash, _ = blake2b.New256(key[0:64])
	}
	s := &Signer{
		hash: hash,
	}
	DefaultFormatter(s)

	// Leave caller options til last so that they override defaults.
	for _, o := range opts {
		o(s)
	}

	return s
}

// Option permits customising the construction of a Signer
type Option func(*Signer)

// SkipQuery instructs Signer to skip the query string when calculating the
// signature. This is useful, say, if you have pagination query parameters but
// you want to use the same signed URL regardless of their value.
func SkipQuery() Option {
	return func(s *Signer) {
		s.skipQuery = true
	}
}

// PrefixPath prefixes the signed URL's path with a string. This can make it easier for a server
// to differentiate between signed and non-signed URLs. Note: the prefix is not
// part of the signature calculation.
func PrefixPath(prefix string) Option {
	return func(s *Signer) {
		s.prefix = prefix
	}
}

// WithQueryFormatter instructs Signer to use query parameters to store the signature
// and expiry in a signed URL.
func WithQueryFormatter() Option {
	return func(s *Signer) {
		s.Formatter = &QueryFormatter{s}
	}
}

// WithPathFormatter instructs Signer to store the signature and expiry in the
// path of a signed URL.
func WithPathFormatter() Option {
	return func(s *Signer) {
		s.Formatter = &PathFormatter{s}
	}
}

// Sign generates a signed URL with the given lifespan.
func (s *Signer) Sign(unsigned string, lifespan time.Duration) (string, error) {
	u, err := url.ParseRequestURI(unsigned)
	if err != nil {
		return "", err
	}

	if s.skipQuery {
		// remove query from signature calculation
		u.RawQuery = ""
	}

	expiry := time.Now().Add(lifespan)
	s.AddExpiry(u, expiry)

	// sign payload creating a signature
	sig := s.sign([]byte(u.String()))

	// add signature to url
	s.AddSignature(u, sig)

	if s.prefix != "" {
		u.Path = path.Join(s.prefix, u.Path)
	}

	// return signed URL
	return u.String(), nil
}

// Verify verifies a signed URL
func (s *Signer) Verify(signed string) error {
	u, err := url.ParseRequestURI(signed)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(u.Path, s.prefix) {
		return ErrInvalidSignedURL
	}
	u.Path = u.Path[len(s.prefix):]

	// Extract signature from url, returning signature and url without
	// signature, which is the input for the signature computation.
	u, sig, err := s.ExtractSignature(u)
	if err != nil {
		return err
	}

	// create another signature for comparison
	compare := s.sign([]byte(u.String()))
	if subtle.ConstantTimeCompare(sig, compare) != 1 {
		return ErrInvalidSignature
	}

	// get expiry from payload
	_, expiry, err := s.ExtractExpiry(u)
	if err != nil {
		return err
	}
	if time.Now().After(expiry) {
		return ErrExpired
	}

	// valid, unexpired, signature
	return nil
}

func (s *Signer) sign(data []byte) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.dirty {
		s.hash.Reset()
	}
	s.dirty = true
	s.hash.Write(data)
	return s.hash.Sum(nil)
}
