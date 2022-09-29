package signer

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	sign := New([]byte("abc123"))

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "with query",
			url:  "https://example.com/a/b/c?baz=cow&foo=bar",
		},
		{
			name: "without query",
			url:  "https://example.com/a/b/c",
		},
		{
			name: "with only question mark",
			url:  "https://example.com/a/b/c?",
		},
		{
			name: "only absolute path",
			url:  "/a/b/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := sign.Sign(tt.url, time.Second*10)
			require.NoError(t, err)

			// check valid URL
			_, err = url.Parse(signed)
			require.NoError(t, err)

			// check valid signature
			err = sign.Verify(signed)
			require.NoError(t, err)
		})
	}

	t.Run("verifiable", func(t *testing.T) {
		// this URL was created with a maximum lifespan:
		// u := "https://example.com/a/b/c?baz=cow&foo=bar"
		// signed, err := sign.Sign(u, time.Duration(math.MaxInt64))
		u := "https://example.com/signed/_BGBJ-6OcP6GnoQz071_rU_VfMWRbi0MGLLQhfxesRg.10887835696/a/b/c?baz=cow&foo=bar"
		err := sign.Verify(u)
		require.NoError(t, err)
	})

	t.Run("expired", func(t *testing.T) {
		u := "https://example.com/a/b/c?baz=cow&foo=bar"
		signed, err := sign.Sign(u, time.Duration(0))
		require.NoError(t, err)

		err = sign.Verify(signed)
		assert.Equal(t, ErrExpired, err)
	})

	t.Run("relative path", func(t *testing.T) {
		_, err := sign.Sign("foo/bar", time.Minute)
		assert.Error(t, err)
	})

	t.Run("invalid prefix", func(t *testing.T) {
		err := sign.Verify("http://abc.com/wrongprefix/fJLFKJ3903.123/foo/bar")
		assert.Equal(t, ErrInvalidMessageFormat, err)
	})

	t.Run("invalid format", func(t *testing.T) {
		err := sign.Verify("http://abc.com/signed/fkljjlFJ903$123/foo/bar")
		assert.Equal(t, ErrInvalidMessageFormat, err)
	})

	t.Run("invalid signature", func(t *testing.T) {
		err := sign.Verify("http://abc.com/signed/MICKEYMOUSE.123/foo/bar")
		assert.Equal(t, ErrInvalidSignature, err)
	})

	t.Run("empty url", func(t *testing.T) {
		_, err := sign.Sign("", 10*time.Second)
		assert.Error(t, err)
	})

	t.Run("not a url", func(t *testing.T) {
		_, err := sign.Sign("cod", 10*time.Second)
		assert.Error(t, err)
	})
}
