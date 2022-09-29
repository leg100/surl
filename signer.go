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

var (
	// ErrInvalidSignature is returned when the provided token's
	// signatuire is not valid.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidMessageFormat is returned when the message's format is
	// invalid.
	ErrInvalidMessageFormat = errors.New("invalid message format")
	// ErrExpired is returned by when the signed URL's expiry has been
	// exceeded.
	ErrExpired = errors.New("URL has expired")
)

// Signer is capable of signing and verifying signed URLs with an expiry.
type Signer struct {
	mu        sync.Mutex
	hash      hash.Hash
	dirty     bool
	skipQuery bool

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
		// This is a case that is so easily avoidable when using this pacakge
		// and since chaining is convenient for this package.  We're going
		// to do the below to handle this possible case so we don't have
		// to return an error.
		hash, _ = blake2b.New256(key[0:64])
	}
	s := &Signer{
		hash:      hash,
		Formatter: &URLPathFormatter{Prefix: "/signed/"},
	}
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

// Sign generates a signed URL with the given lifespan.
func (s *Signer) Sign(u string, lifespan time.Duration) (string, error) {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return "", err
	}
	// retrieve signable part of URL
	data, err := s.parseURL(parsed)
	if err != nil {
		return "", err
	}

	// calculate expiry
	exp := time.Now().Add(lifespan)
	// add expiry to data to create the payload
	payload := s.AddExpiry(exp, []byte(data))
	// sign payload creating a signature
	signature := s.sign(payload)
	// add signature to payload to create the signed data
	data = s.AddSignature(signature, payload)

	// return updated URL
	return s.updateURL(parsed, data), nil
}

// Verify verifies a signed URL
func (s *Signer) Verify(u string) error {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}
	// retrieve signable part of URL
	data, err := s.parseURL(parsed)
	if err != nil {
		return err
	}

	// get signature and payload from data
	signature, payload, err := s.ExtractSignature([]byte(data))
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

// updateURL updates the unsigned url with the signable part.
func (s *Signer) updateURL(u *url.URL, msg []byte) string {
	if s.skipQuery {
		u.Path = string(msg)
		return u.String()
	}

	questionMark := 0
	for i, b := range msg {
		if b == '?' {
			questionMark = i
			break
		}
	}
	if questionMark != 0 {
		u.Path = string(msg[:questionMark])
		// check whether there is anything after '?'
		if len(msg) > questionMark {
			u.RawQuery = string(msg[questionMark+1:])
		}
	} else {
		u.Path = string(msg)
	}
	return u.String()
}

// parseURL parses the signable part of an unsigned URL, which is the path and
// optionally the query as well.
func (s *Signer) parseURL(u *url.URL) ([]byte, error) {
	signable := u.Path
	if !s.skipQuery {
		if u.RawQuery != "" {
			signable = signable + "?" + u.RawQuery
		}
	}
	return []byte(signable), nil
}
