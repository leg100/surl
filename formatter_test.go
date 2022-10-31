package surl

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		data   string
	}{
		{
			name: "without prefix",
			data: "/foo/bar",
		},
		{
			name:   "with prefix",
			data:   "/foo/bar",
			prefix: "/signed/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := URLPathFormatter{tt.prefix}

			exp := time.Date(2081, time.February, 24, 4, 0, 0, 0, time.UTC)

			payload := f.AddExpiry(exp, []byte(tt.data))
			assert.Equal(t, path.Join("3507595200", tt.data), string(payload))

			msg := f.AddSignature([]byte("abcdef"), payload)
			assert.Equal(t, tt.prefix+path.Join("YWJjZGVm.3507595200", tt.data), string(msg))

			sig, payload, err := f.ExtractSignature(msg)
			require.NoError(t, err)
			assert.Equal(t, "abcdef", string(sig))
			assert.Equal(t, path.Join("3507595200", tt.data), string(payload))

			exp, data, err := f.ExtractExpiry(payload)
			require.NoError(t, err)
			assert.Equal(t, exp, exp)
			assert.Equal(t, tt.data, string(data))
		})
	}
}
