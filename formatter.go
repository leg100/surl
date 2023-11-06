package surl

import (
	"net/url"
)

// formatter adds/extracts the signature and expiry to/from a URL according to a
// specific format
type formatter interface {
	// addExpiry adds an expiry to the unsigned URL
	addExpiry(unsigned *url.URL, expiry string)
	// buildPayload produces a payload for signature computation
	buildPayload(url.URL, payloadOptions) string
	// addSignature adds a signature to a URL
	addSignature(*url.URL, string)
	// extractSignature extracts a signature from a URL
	extractSignature(*url.URL) (string, error)
	// extractExpiry extracts an expiry from a URL
	extractExpiry(*url.URL) (string, error)
}

// payloadOptions are options that alter the payload to be signed.
type payloadOptions struct {
	skipQuery  bool
	skipScheme bool
}
