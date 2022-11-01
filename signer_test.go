package surl

import (
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	formatters := []struct {
		name      string
		formatter Option
	}{
		{
			name:      "by_path",
			formatter: WithPathFormatter(),
		},
		{
			name:      "by_query",
			formatter: WithQueryFormatter(),
		},
	}

	opts := []struct {
		name    string
		options []Option
	}{
		{
			name: "no opts",
		},
		{
			name:    "with prefix",
			options: []Option{PrefixPath("/signed")},
		},
		{
			name:    "skip query",
			options: []Option{SkipQuery()},
		},
		{
			name:    "skip query and with prefix",
			options: []Option{SkipQuery(), PrefixPath("/signed")},
		},
	}

	inputs := []struct {
		name     string
		unsigned string
	}{
		{
			name:     "with query",
			unsigned: "https://example.com/a/b/c?foo=bar",
		},
		{
			name:     "without query",
			unsigned: "https://example.com/a/b/c",
		},
		{
			name:     "with only question mark",
			unsigned: "https://example.com/a/b/c?",
		},
		{
			name:     "absolute path",
			unsigned: "/a/b/c",
		},
	}
	// invoke test for each combination of unsigned url, formatter, and set of
	// options
	for _, tt := range inputs {
		for _, f := range formatters {
			for _, opt := range opts {
				options := append(opt.options, f.formatter)
				signer := New([]byte("abc123"), options...)

				t.Run(path.Join(tt.name, f.name, opt.name), func(t *testing.T) {
					signed, err := signer.Sign(tt.unsigned, time.Second*10)
					require.NoError(t, err)

					// check valid URL
					_, err = url.Parse(signed)
					require.NoError(t, err)

					// check valid signature
					err = signer.Verify(signed)
					require.NoError(t, err)
				})
			}
		}
	}
}

func TestSigner_SkipQuery(t *testing.T) {
	signer := New([]byte("abc123"))

	// Demonstrate the SkipQuery option by changing the
	// query string on the signed URL and showing it still verifies.
	t.Run("skip query", func(t *testing.T) {
		sign := New([]byte("abc123"), SkipQuery())

		u := "https://example.com/a/b/c?foo=bar"
		signed, err := sign.Sign(u, time.Minute)
		require.NoError(t, err)

		signed = signed + "&page_num=3&page_size=20"

		err = sign.Verify(signed)
		require.NoError(t, err)
	})

	// Demonstrate how changing the query string invalidates the signed URL
	t.Run("do not skip query", func(t *testing.T) {
		u := "https://example.com/a/b/c?foo=bar"
		signed, err := signer.Sign(u, time.Minute)
		require.NoError(t, err)

		signed = signed + "&page_num=3&page_size=20"

		err = signer.Verify(signed)
		assert.Equal(t, ErrInvalidSignature, err)
	})
}

func TestSigner_Prefix(t *testing.T) {
	signer := New([]byte("abc123"), PrefixPath("/signed"))

	t.Run("invalid prefix", func(t *testing.T) {
		err := signer.Verify("http://abc.com/wrongprefix/foo/bar?expiry=123&signature=fJLFKJ3903")
		assert.Equal(t, ErrInvalidSignedURL, err)
	})
}

func TestSigner_Errors(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		signer := New([]byte("abc123"))

		u := "https://example.com/a/b/c?baz=cow&foo=bar"
		signed, err := signer.Sign(u, time.Duration(0))
		require.NoError(t, err)

		err = signer.Verify(signed)
		assert.Equal(t, ErrExpired, err)
	})

	t.Run("relative path", func(t *testing.T) {
		signer := New([]byte("abc123"))
		_, err := signer.Sign("foo/bar", time.Minute)
		assert.Error(t, err)
	})

	t.Run("empty url", func(t *testing.T) {
		signer := New([]byte("abc123"))
		_, err := signer.Sign("", 10*time.Second)
		assert.Error(t, err)
	})

	t.Run("not a url", func(t *testing.T) {
		signer := New([]byte("abc123"))
		_, err := signer.Sign("cod", 10*time.Second)
		assert.Error(t, err)
	})
}
