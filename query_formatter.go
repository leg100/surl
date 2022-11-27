package surl

import (
	"fmt"
	"net/url"
)

type queryFormatter struct{}

func (f *queryFormatter) addExpiry(unsigned *url.URL, expiry string) {
	q := unsigned.Query()
	q.Add("expiry", expiry)
	unsigned.RawQuery = q.Encode()
}

func (f *queryFormatter) buildPayload(u url.URL, opts payloadOptions) string {
	if opts.skipQuery {
		// Remove all query params other than expiry
		expiry := u.Query().Get("expiry")
		u.RawQuery = url.Values{"expiry": []string{expiry}}.Encode()
	}
	return u.String()
}

func (f *queryFormatter) addSignature(payload *url.URL, sig string) {
	q := payload.Query()
	q.Add("signature", sig)
	payload.RawQuery = q.Encode()
}

func (f *queryFormatter) extractSignature(u *url.URL) (string, error) {
	q := u.Query()
	sig := q.Get("signature")
	if sig == "" {
		return "", fmt.Errorf("%w: %s", ErrInvalidSignedURL, u.String())
	}
	q.Del("signature")
	u.RawQuery = q.Encode()

	return sig, nil
}

func (f *queryFormatter) extractExpiry(u *url.URL) (string, error) {
	q := u.Query()
	expiry := q.Get("expiry")
	if expiry == "" {
		return "", ErrInvalidSignedURL
	}
	q.Del("expiry")
	u.RawQuery = q.Encode()

	return expiry, nil
}
