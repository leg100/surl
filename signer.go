package surl

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
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
	// ErrExpired is returned when the signed URL's expiry has been
	// exceeded.
	ErrExpired = errors.New("URL has expired")

	// Default formatter is the query formatter.
	DefaultFormatter = WithQueryFormatter()
	// Default expiry encoding is base10 (decimal)
	DefaultExpiryFormatter = WithDecimalExpiry()
)

// Signer is capable of signing and verifying signed URLs with an expiry.
type Signer struct {
	mu          sync.Mutex
	hash        hash.Hash
	dirty       bool
	prefix      string
	payloadOpts payloadOptions

	formatter
	intEncoding
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
	DefaultExpiryFormatter(s)

	// Leave caller options til last so that they override defaults.
	for _, o := range opts {
		o(s)
	}

	return s
}

// Option permits customising the construction of a Signer
type Option func(*Signer)

// SkipQuery instructs Signer to skip the query string when computing the
// signature. This is useful, say, if you have pagination query parameters but
// you want to use the same signed URL regardless of their value.
func SkipQuery() Option {
	return func(s *Signer) {
		s.payloadOpts.SkipQuery = true
	}
}

// PrefixPath prefixes the signed URL's path with a string. This can make it easier for a server
// to differentiate between signed and non-signed URLs. Note: the prefix is not
// part of the signature computation.
func PrefixPath(prefix string) Option {
	return func(s *Signer) {
		s.prefix = prefix
	}
}

// WithQueryFormatter instructs Signer to use query parameters to store the signature
// and expiry in a signed URL.
func WithQueryFormatter() Option {
	return func(s *Signer) {
		s.formatter = &queryFormatter{}
	}
}

// WithPathFormatter instructs Signer to store the signature and expiry in the
// path of a signed URL.
func WithPathFormatter() Option {
	return func(s *Signer) {
		s.formatter = &pathFormatter{}
	}
}

// WithDecimalExpiry instructs Signer to use base10 to encode the expiry
func WithDecimalExpiry() Option {
	return func(s *Signer) {
		s.intEncoding = stdIntEncoding(10)
	}
}

// WithBase58Expiry instructs Signer to use base58 to encode the expiry
func WithBase58Expiry() Option {
	return func(s *Signer) {
		s.intEncoding = &base58Encoding{}
	}
}

// Sign generates a signed URL with the given lifespan.
func (s *Signer) Sign(unsigned string, lifespan time.Duration) (string, error) {
	u, err := url.ParseRequestURI(unsigned)
	if err != nil {
		return "", err
	}

	// Add expiry to unsigned URL
	expiry := time.Now().Add(lifespan)
	encodedExpiry := s.Encode(expiry.Unix())
	s.AddExpiry(u, encodedExpiry)

	// Build payload for signature computation
	payload := s.BuildPayload(*u, s.payloadOpts)

	// Sign payload creating a signature
	sig := s.sign([]byte(payload))

	// Add signature to url
	encodedSig := base64.RawURLEncoding.EncodeToString(sig)
	s.AddSignature(u, encodedSig)

	if s.prefix != "" {
		u.Path = path.Join(s.prefix, u.Path)
	}

	// return signed URL
	return u.String(), nil
}

// Verify verifies a signed URL, validating its signature and ensuring it is
// unexpired.
func (s *Signer) Verify(signed string) error {
	u, err := url.ParseRequestURI(signed)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(u.Path, s.prefix) {
		return ErrInvalidSignedURL
	}
	u.Path = u.Path[len(s.prefix):]

	encodedSig, err := s.ExtractSignature(u)
	if err != nil {
		return err
	}
	sig, err := base64.RawURLEncoding.DecodeString(encodedSig)
	if err != nil {
		return fmt.Errorf("%w: invalid base64: %s", ErrInvalidSignature, encodedSig)
	}

	// build the payload for signature computation
	payload := s.BuildPayload(*u, s.payloadOpts)

	// create another signature for comparison and compare
	compare := s.sign([]byte(payload))
	if subtle.ConstantTimeCompare(sig, compare) != 1 {
		return ErrInvalidSignature
	}

	// get expiry from signed URL
	encodedExpiry, err := s.ExtractExpiry(u)
	if err != nil {
		return err
	}
	expiry, err := s.Decode(encodedExpiry)
	if err != nil {
		return err
	}
	if time.Now().After(time.Unix(expiry, 0)) {
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
