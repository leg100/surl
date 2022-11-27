package surl

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathFormatter(t *testing.T) {
	f := pathFormatter{}
	// unsigned url with existing path
	u := &url.URL{Path: "/foo/bar"}

	expiry := time.Date(2081, time.February, 24, 4, 0, 0, 0, time.UTC)
	encoded := stdIntEncoding(10).Encode(expiry.Unix())

	f.AddExpiry(u, encoded)
	assert.Equal(t, "3507595200/foo/bar", u.Path)

	f.AddSignature(u, "abcdef")
	assert.Equal(t, "/abcdef.3507595200/foo/bar", u.Path)

	sig, err := f.ExtractSignature(u)
	require.NoError(t, err)
	assert.Equal(t, "abcdef", string(sig))
	assert.Equal(t, "3507595200/foo/bar", u.Path)

	got, err := f.ExtractExpiry(u)
	require.NoError(t, err)
	assert.Equal(t, encoded, got)
	assert.Equal(t, "/foo/bar", u.Path)
}

func TestPathFormatter_Errors(t *testing.T) {
	signer := New([]byte("abc123"), WithPathFormatter())

	t.Run("invalid signature", func(t *testing.T) {
		err := signer.Verify("http://abc.com/MICKEYMOUSE.123/foo/bar")
		assert.Truef(t, errors.Is(err, ErrInvalidSignature), "got error: %w", err)
	})
}
