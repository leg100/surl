package surl

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryFormatter(t *testing.T) {
	f := QueryFormatter{&Signer{}}
	// unsigned url with existing query
	u := &url.URL{RawQuery: "foo=bar"}

	expiry := time.Date(2081, time.February, 24, 4, 0, 0, 0, time.UTC)
	encoded := StdIntEncoding(10).Encode(expiry.Unix())

	f.AddExpiry(u, encoded)
	assert.Equal(t, "expiry=3507595200&foo=bar", u.RawQuery)

	f.AddSignature(u, []byte("abcdef"))
	assert.Equal(t, "expiry=3507595200&foo=bar&signature=YWJjZGVm", u.RawQuery)

	sig, err := f.ExtractSignature(u)
	require.NoError(t, err)
	assert.Equal(t, "abcdef", string(sig))
	assert.Equal(t, "expiry=3507595200&foo=bar", u.RawQuery)

	got, err := f.ExtractExpiry(u)
	require.NoError(t, err)
	assert.Equal(t, encoded, got)
	assert.Equal(t, "foo=bar", u.RawQuery)
}

func TestQueryFormatter_Errors(t *testing.T) {
	signer := New([]byte("abc123"), WithQueryFormatter())

	t.Run("missing query params", func(t *testing.T) {
		err := signer.Verify("https://example.com/a/b/c?foo=bar")
		assert.Truef(t, errors.Is(err, ErrInvalidSignedURL), "got error: %w", err)
	})

	t.Run("missing signature param", func(t *testing.T) {
		err := signer.Verify("https://example.com/a/b/c?foo=bar&expiry=123")
		assert.True(t, errors.Is(err, ErrInvalidSignedURL), "got error: %w", err)
	})
}
