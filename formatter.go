package surl

import (
	"net/url"
)

// Formatter adds/extracts the signature and expiry to/from a URL according to a
// specific format
type Formatter interface {
	// AddExpiry adds an expiry to the unsigned URL
	AddExpiry(unsigned *url.URL, expiry string)
	// BuildPayload produces a payload for signature computation
	BuildPayload(url.URL, PayloadOptions) string
	// AddSignature adds a signature to a URL
	AddSignature(*url.URL, string)
	// ExtractSignature extracts a signature from a URL
	ExtractSignature(*url.URL) (string, error)
	// ExtractExpiry extracts an expiry from a URL
	ExtractExpiry(*url.URL) (string, error)
}

type PayloadOptions struct {
	SkipQuery bool
}
