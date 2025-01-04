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
	// ErrInvalidSignature is returned when the signature is invalid.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidFormat is returned when the format of the signed URL is
	// invalid.
	ErrInvalidFormat = errors.New("invalid format")
	// ErrExpired is returned when a signed URL has expired.
	ErrExpired = errors.New("URL has expired")

	// Default formatter is the query formatter.
	DefaultFormatter = WithQueryFormatter()
	// Default expiry encoding is base10 (decimal)
	DefaultExpiryFormatter = WithDecimalExpiry()
)

// Signer is capable of signing and verifying signed URLs with an expiry.
type Signer struct {
	mu     sync.Mutex
	hash   hash.Hash
	dirty  bool
	prefix string

	payloadOptions
	formatter
	intEncoding
}

// New constructs a new signer, performing the one-off task of generating a
// secure hash from the key. The key must be between 0 and 64 bytes long;
// anything longer is truncated. Options alter the default format and behaviour
// of signed URLs.
func New(key []byte, opts ...Option) *Signer {
	hash, err := blake2b.New256(key)
	if err != nil {
		// Safely ignore one and only error regarding keys longer than 64 bytes.
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
		s.skipQuery = true
	}
}

// SkipScheme instructs Signer to skip the scheme when computing the signature.
// This is useful, say, if you generate signed URLs in production where you use
// https but you want to use these URLs in development too where you use http.
func SkipScheme() Option {
	return func(s *Signer) {
		s.skipScheme = true
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
func (s *Signer) Sign(unsigned string, expiry time.Time) (string, error) {
	u, err := url.ParseRequestURI(unsigned)
	if err != nil {
		return "", err
	}

	// Add expiry to unsigned URL
	encodedExpiry := s.Encode(expiry.Unix())
	s.addExpiry(u, encodedExpiry)

	// Build payload for signature computation
	payload := s.buildPayload(*u, s.payloadOptions)

	// Sign payload creating a signature
	sig := s.sign([]byte(payload))

	// Add signature to url
	encodedSig := base64.RawURLEncoding.EncodeToString(sig)
	s.addSignature(u, encodedSig)

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
		return ErrInvalidFormat
	}
	u.Path = u.Path[len(s.prefix):]

	encodedSig, err := s.extractSignature(u)
	if err != nil {
		return err
	}
	sig, err := base64.RawURLEncoding.DecodeString(encodedSig)
	if err != nil {
		return fmt.Errorf("%w: invalid base64: %s", ErrInvalidSignature, encodedSig)
	}

	// build the payload for signature computation
	payload := s.buildPayload(*u, s.payloadOptions)

	// create another signature for comparison and compare
	compare := s.sign([]byte(payload))
	if subtle.ConstantTimeCompare(sig, compare) != 1 {
		return ErrInvalidSignature
	}

	// get expiry from signed URL
	encodedExpiry, err := s.extractExpiry(u)
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
