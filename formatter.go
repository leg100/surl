package signer

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"time"
)

// Formatter adds/extracts the signature and expiry to/from a URL according to a
// specific format
type Formatter interface {
	// AddExpiry adds the expiry to the data, creating a payload for signing
	AddExpiry(exp time.Time, data []byte) []byte
	// AddSignature adds the signature to the payload, creating a signed message
	AddSignature(sig, payload []byte) []byte
	// ExtractSignature extracts the signature from the signed message,
	// returning the signature as well as the signed payload.
	ExtractSignature(msg []byte) ([]byte, []byte, error)
	// ExtractExpiry extracts the expiry from the signed payload, returning the
	// expiry as well as the original data.
	ExtractExpiry(payload []byte) (time.Time, []byte, error)
}

// URLPathFormatter includes the signature and expiry in a
// message according to the format: <prefix><sig>.<exp>/<data>. Suitable for
// URL paths as an alternative to using query parameters.
type URLPathFormatter struct {
	// Prefix message with a string
	//
	// TODO: default to '/'?
	Prefix string
}

// AddExpiry adds expiry as a base64 encoded component e.g. /foo/bar ->
// 390830893/foo/bar
func (u *URLPathFormatter) AddExpiry(exp time.Time, data []byte) []byte {
	// convert expiry to bytes
	expBytes := strconv.FormatInt(exp.Unix(), 10)

	payload := make([]byte, 0, len(expBytes)+len(data))
	// add expiry
	payload = append(payload, expBytes...)
	// add data
	return append(payload, data...)
}

// AddSignature adds signature as a path component alongside the expiry e.g.
// abZ3G/foo/bar -> KKLJjd3090fklaJKLJK.abZ3G/foo/bar
func (u *URLPathFormatter) AddSignature(sig, payload []byte) []byte {
	encSize := base64.RawURLEncoding.EncodedLen(len(sig))
	// calculate msg capacity
	mcap := encSize + len(payload) + 1 // +1 for '.'
	if u.Prefix != "" {
		mcap += len(u.Prefix)
	}
	msg := make([]byte, 0, mcap)
	// add prefix
	msg = append(msg, u.Prefix...)
	// add encoded sig
	msg = msg[0 : len(u.Prefix)+encSize]
	base64.RawURLEncoding.Encode(msg[len(u.Prefix):], sig)
	// add '.'
	msg = append(msg, '.')
	// add payload
	return append(msg, payload...)
}

// ExtractSignature decodes and splits the signature and payload from the signed message.
func (u *URLPathFormatter) ExtractSignature(msg []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(msg, []byte(u.Prefix)) {
		return nil, nil, ErrInvalidMessageFormat
	}
	// remove prefix
	msg = msg[len(u.Prefix):]

	// prise apart sig and payload
	parts := bytes.SplitN(msg, []byte{'.'}, 2)
	if len(parts) != 2 {
		return nil, nil, ErrInvalidMessageFormat
	}
	sig := parts[0]
	payload := parts[1]

	// decode base64-encoded sig into bytes
	decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(sig)))
	_, err := base64.RawURLEncoding.Decode(decoded, sig)
	if err != nil {
		return nil, nil, err
	}

	return decoded, payload, nil
}

// ExtractExpiry decodes and splits the expiry and data from the payload.
func (u *URLPathFormatter) ExtractExpiry(payload []byte) (time.Time, []byte, error) {
	// prise apart expiry and data
	slash := 0
	for i, b := range payload {
		if b == '/' {
			slash = i
			break
		}
	}
	if slash == 0 {
		return time.Time{}, nil, ErrInvalidMessageFormat
	}
	expBytes := payload[:slash]
	data := payload[slash:]

	// convert bytes into int
	expInt, err := strconv.ParseInt(string(expBytes), 10, 64)
	if err != nil {
		return time.Time{}, nil, err
	}
	// convert int into time.Time
	t := time.Unix(expInt, 0)

	return t, data, nil
}
