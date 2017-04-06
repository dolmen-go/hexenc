// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hexenc provides the same interface as stdlib encoding/base64 for encoding/hex.
*/
package hexenc

import (
	"encoding/hex"
)

// An Encoding is an hexadecimal encoding/decoding scheme.
type Encoding struct{}

// Decode decodes src into DecodedLen(len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// Decode expects that src contain only hexadecimal
// characters and that src should have an even length.
func (Encoding) Decode(dst, src []byte) (n int, err error) {
	return hex.Decode(dst, src)
}

// DecodeString returns the bytes represented by the hexadecimal string s.
func (Encoding) DecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 2.
func (Encoding) DecodedLen(x int) int {
	return hex.DecodedLen(x)
}

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements hexadecimal encoding.
func (Encoding) Encode(dst, src []byte) {
	hex.Encode(dst, src)
}

// EncodeToString returns the hexadecimal encoding of src.
func (Encoding) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 2.
func (Encoding) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}
