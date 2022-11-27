package surl

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"

	"github.com/itchyny/base58-go"
)

// intEncoding encodes integers into strings, and decodes strings into integers
type intEncoding interface {
	Encode(int64) string
	Decode(string) (int64, error)
}

type stdIntEncoding int

func (b stdIntEncoding) Encode(i int64) string {
	return strconv.FormatInt(i, int(b))
}

func (b stdIntEncoding) Decode(s string) (int64, error) {
	return strconv.ParseInt(s, int(b), 64)
}

type base58Encoding struct{}

func (base58Encoding) Encode(i int64) string {
	return string(base58.FlickrEncoding.EncodeUint64(uint64(i)))
}

func (base58Encoding) Decode(s string) (int64, error) {
	i, err := base58.FlickrEncoding.DecodeUint64([]byte(s))
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

type base64Encoding struct{}

func (base64Encoding) Encode(i int64) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return base64.RawURLEncoding.EncodeToString(b)
}

func (base64Encoding) Decode(s string) (int64, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(bytes)), nil
}
