package surl

import (
	"net/url"
	"strings"
)

type pathFormatter struct{}

func (f *pathFormatter) addExpiry(unsigned *url.URL, expiry string) {
	unsigned.Path = expiry + unsigned.Path
}

func (f *pathFormatter) buildPayload(u url.URL, opts payloadOptions) string {
	if opts.skipQuery {
		u.RawQuery = ""
	}
	if opts.skipScheme {
		u.Scheme = ""
	}
	return u.String()
}

func (f *pathFormatter) addSignature(payload *url.URL, sig string) {
	payload.Path = "/" + sig + "." + payload.Path
}

func (f *pathFormatter) extractSignature(u *url.URL) (string, error) {
	// prise apart sig and payload
	sig, payload, found := strings.Cut(u.Path, ".")
	if !found {
		return "", ErrInvalidFormat
	}
	// remove leading /
	sig = sig[1:]

	u.Path = payload

	return sig, nil
}

func (*pathFormatter) extractExpiry(u *url.URL) (string, error) {
	// prise apart expiry and data
	expiry, path, found := strings.Cut(u.Path, "/")
	if !found {
		return "", ErrInvalidFormat
	}
	// add leading slash back to path
	u.Path = "/" + path

	return expiry, nil
}
