package surl

import (
	"fmt"
	"net/url"
)

// QueryFormatter includes the signature and expiry as URL query parameters
// according to the format: /path?expiry=<exp>&signature=<sig>.
type QueryFormatter struct{}

func (f *QueryFormatter) AddExpiry(unsigned *url.URL, expiry string) {
	q := unsigned.Query()
	q.Add("expiry", expiry)
	unsigned.RawQuery = q.Encode()
}

func (f *QueryFormatter) BuildPayload(u url.URL, opts PayloadOptions) string {
	if opts.SkipQuery {
		// Remove all query params other than expiry
		expiry := u.Query().Get("expiry")
		u.RawQuery = url.Values{"expiry": []string{expiry}}.Encode()
	}
	return u.String()
}

// AddSignature adds signature as a query parameter alongside the expiry e.g.
// /foo/bar?expiry=<exp> -> /foo/bar?expiry=<exp>&signature=<sig>
func (f *QueryFormatter) AddSignature(payload *url.URL, sig string) {
	q := payload.Query()
	q.Add("signature", sig)
	payload.RawQuery = q.Encode()
}

// ExtractSignature splits the signature and payload from the signed message.
func (f *QueryFormatter) ExtractSignature(u *url.URL) (string, error) {
	q := u.Query()
	sig := q.Get("signature")
	if sig == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidSignedURL, u.String())
	}
	q.Del("signature")
	u.RawQuery = q.Encode()

	return sig, nil
}

// ExtractExpiry splits the expiry and data from the payload.
func (f *QueryFormatter) ExtractExpiry(u *url.URL) (string, error) {
	q := u.Query()
	expiry := q.Get("expiry")
	if expiry == "" {
		return "", ErrInvalidSignedURL
	}
	q.Del("expiry")
	u.RawQuery = q.Encode()

	return expiry, nil
}
