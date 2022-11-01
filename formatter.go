package surl

import (
	"net/url"
	"time"
)

// Formatter adds/extracts the signature and expiry to/from a URL according to a
// specific format
type Formatter interface {
	// AddExpiry adds an expiry to a URL
	AddExpiry(*url.URL, time.Time)

	// AddSignature adds a signature to a URL
	AddSignature(*url.URL, []byte)

	// ExtractSignature extracts a signature from a URL, returning the modified
	// URL and the signature.
	ExtractSignature(*url.URL) (*url.URL, []byte, error)

	// ExtractExpiry extracts an expiry from a URL, returning the modified URL
	// and the signature.
	ExtractExpiry(*url.URL) (*url.URL, time.Time, error)
}
