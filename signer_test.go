package surl

import (
	"crypto/rand"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	formatters = []struct {
		name      string
		formatter Option
	}{
		{
			name:      "path",
			formatter: WithPathFormatter(),
		},
		{
			name:      "query",
			formatter: WithQueryFormatter(),
		},
	}

	encoders = []struct {
		name    string
		encoder Option
	}{
		{
			name:    "decimal",
			encoder: WithDecimalExpiry(),
		},
		{
			name:    "base58",
			encoder: WithBase58Expiry(),
		},
	}

	opts = []struct {
		name    string
		options []Option
	}{
		{
			name: "no opts",
		},
		{
			name:    "prefix",
			options: []Option{PrefixPath("/signed")},
		},
		{
			name:    "skip query",
			options: []Option{SkipQuery()},
		},
		{
			name:    "skip scheme",
			options: []Option{SkipScheme()},
		},
		{
			name:    "prefix and skip query",
			options: []Option{SkipQuery(), PrefixPath("/signed")},
		},
		{
			name:    "prefix and skip scheme",
			options: []Option{SkipScheme(), PrefixPath("/signed")},
		},
		{
			name:    "prefix and skip query and skip scheme",
			options: []Option{SkipQuery(), SkipScheme(), PrefixPath("/signed")},
		},
	}
)

func TestSigner(t *testing.T) {
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
	// invoke test for each combination of unsigned url, formatter, encoder, and set of
	// options
	for _, tt := range inputs {
		for _, f := range formatters {
			for _, enc := range encoders {
				for _, opt := range opts {
					options := append(opt.options, f.formatter, enc.encoder)
					signer := New([]byte("abc123"), options...)

					t.Run(path.Join(tt.name, f.name, enc.name, opt.name), func(t *testing.T) {
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

		t.Run("check original query parameters are intact", func(t *testing.T) {
			u, err := url.Parse(signed)
			require.NoError(t, err)
			assert.Equal(t, "bar", u.Query().Get("foo"))
		})
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

func TestSigner_SkipScheme(t *testing.T) {
	// Demonstrate the SkipScheme option by changing the scheme on the signed
	// URL and showing it still verifies.
	t.Run("skip scheme", func(t *testing.T) {
		signer := New([]byte("abc123"), SkipScheme())

		unsigned := "https://example.com/a/b/c?foo=bar"
		signed, err := signer.Sign(unsigned, time.Minute)
		require.NoError(t, err)

		u, err := url.Parse(signed)
		require.NoError(t, err)
		u.Scheme = "http"

		err = signer.Verify(u.String())
		require.NoError(t, err)
	})

	// Demonstrate how changing the scheme invalidates the signed URL
	t.Run("do not skip scheme", func(t *testing.T) {
		signer := New([]byte("abc123"))

		unsigned := "https://example.com/a/b/c?foo=bar"
		signed, err := signer.Sign(unsigned, time.Minute)
		require.NoError(t, err)

		u, err := url.Parse(signed)
		require.NoError(t, err)
		u.Scheme = "http"

		err = signer.Verify(u.String())
		assert.Equal(t, ErrInvalidSignature, err)
	})
}

func TestSigner_Prefix(t *testing.T) {
	signer := New([]byte("abc123"), PrefixPath("/signed"))

	t.Run("invalid prefix", func(t *testing.T) {
		err := signer.Verify("http://abc.com/wrongprefix/foo/bar?expiry=123&signature=fJLFKJ3903")
		assert.Equal(t, ErrInvalidFormat, err)
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

	t.Run("scheme has changed", func(t *testing.T) {
		signer := New([]byte("abc123"))
		signed, err := signer.Sign("https://example.com/a/b/c?baz=cow&foo=bar", 10*time.Second)
		require.NoError(t, err)

		hacked, err := url.Parse(signed)
		require.NoError(t, err)
		hacked.Scheme = "http"

		err = signer.Verify(hacked.String())
		assert.Error(t, err)
	})

	t.Run("hostname has changed", func(t *testing.T) {
		signer := New([]byte("abc123"))
		signed, err := signer.Sign("https://example.com/a/b/c?baz=cow&foo=bar", 10*time.Second)
		require.NoError(t, err)

		hacked, err := url.Parse(signed)
		require.NoError(t, err)
		hacked.Host = "hacked.com:1337"

		err = signer.Verify(hacked.String())
		assert.Error(t, err)
	})
}

var (
	bu   string
	berr error
)

func Benchmark(b *testing.B) {
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	require.NoError(b, err)

	// invoke bench for each combination of formatter, encoder, and set of
	// options
	for _, f := range formatters {
		for _, enc := range encoders {
			for _, opt := range opts {
				options := append(opt.options, f.formatter, enc.encoder)

				b.Run(path.Join("sign", f.name, enc.name, opt.name), func(b *testing.B) {
					signer := New(secret, options...)

					var u string
					for n := 0; n < b.N; n++ {
						// store result to prevent compiler eliminating func call
						u, _ = signer.Sign("https://example.com/a/b/c?x=1&y=2&z=3", time.Hour)
					}
					// store result in pkg var to to prevent compiler eliminating benchmark
					bu = u
				})

				b.Run(path.Join("verify", f.name, enc.name, opt.name), func(b *testing.B) {
					signer := New(secret, options...)
					signed, _ := signer.Sign("https://example.com/a/b/c?x=1&y=2&z=3", time.Hour)

					for n := 0; n < b.N; n++ {
						// store result to prevent compiler eliminating func call
						err = signer.Verify(signed)
					}
					// store result in pkg var to to prevent compiler eliminating benchmark
					berr = err
				})
			}
		}
	}
}
