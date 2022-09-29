package signer

import (
	"crypto/subtle"
	"errors"
	"hash"
	"net/url"
	"sync"
	"time"

	"golang.org/x/crypto/blake2b"
)

// ErrInvalidSignature is returned when the provided token's
// signatuire is not valid.
var ErrInvalidSignature = errors.New("invalid signature")

// ErrInvalidMessageFormat is returned when the message's format is
// invalid.
var ErrInvalidMessageFormat = errors.New("invalid message format")

// ErrExpired is returned by when the signed URL's expiry has been
// exceeded.
var ErrExpired = errors.New("URL has expired")

// Signer is the type for the package. Secret is the signer secret, a lengthy
// and hard to guess string we use to sign things. The secret must not exceed 64 characters.
type Signer struct {
	mu    sync.Mutex
	hash  hash.Hash
	dirty bool

	Formatter
}

// New constructs a new signer, performing the one-off task of generating a
// secure hash from the key. The key must be between 0 and 64 bytes long;
// anything longer is stripped off.
func New(key []byte) *Signer {
	hash, err := blake2b.New256(key)
	if err != nil {
		// The only possible error that can be returned here is if the key
		// is larger than 64 bytes - which the blake2b hash will not accept.
		// This is a case that is so easily avoidable when using this pacakge
		// and since chaining is convenient for this package.  We're going
		// to do the below to handle this possible case so we don't have
		// to return an error.
		hash, _ = blake2b.New256(key[0:64])
	}
	return &Signer{
		hash:      hash,
		Formatter: &URLPathFormatter{Prefix: "/signed/"},
	}
}

// Sign generates a signed URL with the given lifespan.
func (s *Signer) Sign(u string, lifespan time.Duration) (string, error) {
	// verify this is a valid url
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return "", err
	}
	// only the URL path is signed; the remainder of the URL is ignored
	data := parsed.Path
	// calculate expiry
	exp := time.Now().Add(lifespan)
	// add expiry to data to create the payload
	payload := s.AddExpiry(exp, []byte(data))
	// sign payload creating a signature
	signature := s.sign(payload)
	// add signature to payload to create the new path
	path := s.AddSignature(signature, payload)

	// return URL updated with signed path
	parsed.Path = string(path)
	return parsed.String(), nil
}

// Verify verifies a signed URL
func (s *Signer) Verify(u string) error {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}
	// only the path is signed
	path := parsed.Path

	// get signature and payload from path
	signature, payload, err := s.ExtractSignature([]byte(path))
	if err != nil {
		return err
	}
	// create another signature for comparison
	compare := s.sign(payload)
	if subtle.ConstantTimeCompare([]byte(signature), compare) != 1 {
		return ErrInvalidSignature
	}

	// get expiry from payload
	exp, _, err := s.ExtractExpiry(payload)
	if err != nil {
		return err
	}
	if time.Now().After(exp) {
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
