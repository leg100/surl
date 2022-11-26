package surl

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

// QueryFormatter includes the signature and expiry as URL query parameters
// according to the format: /path?expiry=<exp>&signature=<sig>.
type QueryFormatter struct {
	signer *Signer
}

// AddExpiry adds expiry as a query parameter e.g. /foo/bar ->
// /foo/bar?expiry=<exp>
func (f *QueryFormatter) AddExpiry(unsigned *url.URL, expiry string) {
	q := unsigned.Query()
	q.Add("expiry", expiry)
	unsigned.RawQuery = q.Encode()
}

// AddSignature adds signature as a query parameter alongside the expiry e.g.
// /foo/bar?expiry=<exp> -> /foo/bar?expiry=<exp>&signature=<sig>
func (f *QueryFormatter) AddSignature(payload *url.URL, sig []byte) {
	encoded := base64.RawURLEncoding.EncodeToString(sig)

	q := payload.Query()
	q.Add("signature", encoded)
	payload.RawQuery = q.Encode()
}

// ExtractSignature decodes and splits the signature and payload from the signed message.
func (f *QueryFormatter) ExtractSignature(u *url.URL) ([]byte, error) {
	q := u.Query()
	encoded := q.Get("signature")
	if encoded == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSignedURL, u.String())
	}

	// decode base64-encoded sig into bytes
	sig, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if f.signer.skipQuery {
		// remove all query params other than expiry because they don't form
		// part of the input to the signature computation.
		expiry := u.Query().Get("expiry")
		u.RawQuery = url.Values{"expiry": {expiry}}.Encode()
	} else {
		q.Del("signature")
		u.RawQuery = q.Encode()
	}

	return sig, nil
}

// ExtractExpiry decodes and splits the expiry and data from the payload.
func (f *QueryFormatter) ExtractExpiry(u *url.URL) (string, error) {
	q := u.Query()
	encoded := q.Get("expiry")
	if encoded == "" {
		return "", ErrInvalidSignedURL
	}
	q.Del("expiry")
	u.RawQuery = q.Encode()

	return encoded, nil
}
