//go:build !go1.10 || builtinencoder
// +build !go1.10 builtinencoder

package hexenc

import (
	"encoding/hex"
	"io"
)

type encoder struct {
	err error
	w   io.Writer
	out [1024]byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.err != nil {
		return 0, e.err
	}

	for len(p) > 0 {
		nn := len(e.out) / 2
		if nn > len(p) {
			nn = len(p)
		}
		hex.Encode(e.out[:], p[:nn])
		var m int
		if m, e.err = e.w.Write(e.out[:2*nn]); e.err != nil {
			return n + m/2, e.err
		}
		n += nn
		p = p[nn:]
	}
	return
}

// NewEncoder returns an io.Writer that writes lowercase hexadecimal characters
// to w.
func (Encoding) NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

type decoder struct {
	err error
	r   io.Reader
	buf [1024]byte
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	var remaining int
	var b byte

	for len(p) > 0 {
		buf := p[:len(p)&^1]
		if 2*len(p) <= len(d.buf) {
			buf = d.buf[:2*len(p)]
		}
		rbuf := buf
		if remaining > 0 {
			buf[0] = b
			rbuf = buf[1:]
		}
		nn, readErr := d.r.Read(rbuf)
		nn += remaining

		remaining = nn & 1
		if remaining == 1 {
			b = buf[nn-1]
			nn = nn &^ 1
		}

		if nn > 1 {
			var m int
			m, d.err = hex.Decode(p, buf[:nn])
			n += m
			if d.err != nil {
				return n, d.err
			}
			p = p[m:]
		}
		if readErr == io.EOF && remaining == 1 {
			readErr = io.ErrUnexpectedEOF
		}
		d.err = readErr
		if readErr != nil {
			return n, readErr
		}
	}

	return n, nil
}

// NewDecoder returns an io.Reader that decodes hexadecimal characters from r.
// NewDecoder expects that r contain only an even number of hexadecimal
// characters.
func (Encoding) NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}
