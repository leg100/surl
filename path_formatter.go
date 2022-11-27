package surl

import (
	"net/url"
	"strings"
)

// pathFormatter includes the signature and expiry in a
// message according to the format: <sig>.<exp>/<data>. Suitable for
// URL paths as an alternative to using query parameters.
type pathFormatter struct{}

// AddExpiry adds expiry as a path component e.g. /foo/bar ->
// 390830893/foo/bar
func (f *pathFormatter) AddExpiry(unsigned *url.URL, expiry string) {
	unsigned.Path = expiry + unsigned.Path
}

// AddExpiry adds expiry as a path component e.g. /foo/bar ->
// 390830893/foo/bar
func (f *pathFormatter) BuildPayload(u url.URL, opts payloadOptions) string {
	if opts.SkipQuery {
		u.RawQuery = ""
	}
	return u.String()
}

// AddSignature adds signature as a path component alongside the expiry e.g.
// abZ3G/foo/bar -> /KKLJjd3090fklaJKLJK.abZ3G/foo/bar
func (f *pathFormatter) AddSignature(payload *url.URL, sig string) {
	payload.Path = "/" + sig + "." + payload.Path
}

// ExtractSignature splits the signature and payload from the signed URL.
func (f *pathFormatter) ExtractSignature(u *url.URL) (string, error) {
	// prise apart sig and payload
	sig, payload, found := strings.Cut(u.Path, ".")
	if !found {
		return "", ErrInvalidSignedURL
	}
	// remove leading /
	sig = sig[1:]

	u.Path = payload

	return sig, nil
}

// ExtractExpiry splits the expiry and data from the payload.
func (*pathFormatter) ExtractExpiry(u *url.URL) (string, error) {
	// prise apart expiry and data
	expiry, path, found := strings.Cut(u.Path, "/")
	if !found {
		return "", ErrInvalidSignedURL
	}
	// add leading slash back to path
	u.Path = "/" + path

	return expiry, nil
}
