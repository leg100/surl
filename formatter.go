package surl

import (
	"net/url"
)

// Formatter adds/extracts the signature and expiry to/from a URL according to a
// specific format
type Formatter interface {
	// AddExpiry adds an expiry to a URL
	AddExpiry(u *url.URL, expiry string)
	// AddSignature adds a signature to a URL
	AddSignature(*url.URL, []byte)
	// ExtractSignature extracts a signature from a URL
	ExtractSignature(*url.URL) ([]byte, error)
	// ExtractExpiry extracts an expiry from a URL
	ExtractExpiry(*url.URL) (string, error)
}
