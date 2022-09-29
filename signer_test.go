package signer

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignature_SignURL(t *testing.T) {
	sign := New([]byte("abc123"))

	t.Run("signable", func(t *testing.T) {
		u := "https://example.com/a/b/c?baz=cow&foo=bar"
		signed, err := sign.Sign(u, time.Second*10)
		require.NoError(t, err)

		// check valid URL
		_, err = url.Parse(signed)
		require.NoError(t, err)

		// check valid signature
		err = sign.Verify(signed)
		require.NoError(t, err)
	})

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
