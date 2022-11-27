package surl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntEncoding(t *testing.T) {
	tests := []struct {
		name     string
		encoding intEncoding
		input    int64
		want     string
	}{
		{
			name:     "decimal",
			encoding: stdIntEncoding(10),
			input:    3507595200,
			want:     "3507595200",
		},
		{
			name:     "hex",
			encoding: stdIntEncoding(16),
			input:    3507595200,
			want:     "d111a7c0",
		},
		{
			name:     "base58",
			encoding: &base58Encoding{},
			input:    3507595200,
			want:     "6kXkRj",
		},
		{
			name:     "base64",
			encoding: &base64Encoding{},
			input:    3507595200,
			want:     "AAAAANERp8A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.encoding.Encode(tt.input))

			decoded, err := tt.encoding.Decode(tt.want)
			require.NoError(t, err)
			assert.Equal(t, tt.input, decoded)
		})
	}
}
