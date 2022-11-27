package surl

import (
	"net/url"
	"strings"
)

// PathFormatter includes the signature and expiry in a
// message according to the format: <sig>.<exp>/<data>. Suitable for
// URL paths as an alternative to using query parameters.
type PathFormatter struct{}

// AddExpiry adds expiry as a path component e.g. /foo/bar ->
// 390830893/foo/bar
func (f *PathFormatter) AddExpiry(unsigned *url.URL, expiry string) {
	unsigned.Path = expiry + unsigned.Path
}

// AddExpiry adds expiry as a path component e.g. /foo/bar ->
// 390830893/foo/bar
func (f *PathFormatter) BuildPayload(u url.URL, opts PayloadOptions) string {
	if opts.SkipQuery {
		u.RawQuery = ""
	}
	return u.String()
}

// AddSignature adds signature as a path component alongside the expiry e.g.
// abZ3G/foo/bar -> /KKLJjd3090fklaJKLJK.abZ3G/foo/bar
func (f *PathFormatter) AddSignature(payload *url.URL, sig string) {
	payload.Path = "/" + sig + "." + payload.Path
}

// ExtractSignature splits the signature and payload from the signed URL.
func (f *PathFormatter) ExtractSignature(u *url.URL) (string, error) {
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
func (*PathFormatter) ExtractExpiry(u *url.URL) (string, error) {
	// prise apart expiry and data
	expiry, path, found := strings.Cut(u.Path, "/")
	if !found {
		return "", ErrInvalidSignedURL
	}
	// add leading slash back to path
	u.Path = "/" + path

	return expiry, nil
}
