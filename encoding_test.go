// base64.RawStdEncoding appears with go 1.5
//
//go:build go1.5
// +build go1.5

package hexenc

import (
	"encoding/base64"
)

var _ = []interface {
	Decode(dst, src []byte) (n int, err error)
	DecodeString(s string) ([]byte, error)
	DecodedLen(x int) int
	Encode(dst, src []byte)
	EncodeToString(src []byte) string
	EncodedLen(n int) int
}{
	base64.RawStdEncoding,
	base64.RawURLEncoding,
	Encoding{},
}
