//+build go1.10,!builtinencoder

package hexenc

import (
	"encoding/hex"
	"io"
)

// NewEncoder returns an io.Writer that writes lowercase hexadecimal characters
// to w.
func (Encoding) NewEncoder(w io.Writer) io.Writer {
	return hex.NewEncoder(w)
}

// NewDecoder returns an io.Reader that decodes hexadecimal characters from r.
// NewDecoder expects that r contain only an even number of hexadecimal
// characters.
func (Encoding) NewDecoder(r io.Reader) io.Reader {
	return hex.NewDecoder(r)
}
