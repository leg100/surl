package surl

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// QueryFormatter includes the signature and expiry as URL query parameters
// according to the format: /path?expiry=<exp>&signature=<sig>.
type QueryFormatter struct {
	signer *Signer
}

// AddExpiry adds expiry as a query parameter e.g. /foo/bar ->
// /foo/bar?expiry=<exp>
func (f *QueryFormatter) AddExpiry(unsigned *url.URL, exp time.Time) {
	// convert expiry to string
	val := strconv.FormatInt(exp.Unix(), 10)

	q := unsigned.Query()
	q.Add("expiry", val)
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
func (f *QueryFormatter) ExtractSignature(u *url.URL) (*url.URL, []byte, error) {
	q := u.Query()
	encoded := q.Get("signature")
	if encoded == "" {
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidSignedURL, u.String())
	}

	// decode base64-encoded sig into bytes
	sig, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, nil, err
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

	return u, sig, nil
}

// ExtractExpiry decodes and splits the expiry and data from the payload.
func (f *QueryFormatter) ExtractExpiry(u *url.URL) (*url.URL, time.Time, error) {
	q := u.Query()
	encoded := q.Get("expiry")
	if encoded == "" {
		return nil, time.Time{}, ErrInvalidSignedURL
	}
	q.Del("expiry")
	u.RawQuery = q.Encode()

	// convert bytes into int
	expInt, err := strconv.ParseInt(string(encoded), 10, 64)
	if err != nil {
		return nil, time.Time{}, err
	}
	// convert int into time.Time
	t := time.Unix(expInt, 0)

	return u, t, nil
}
