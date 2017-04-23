package hexenc_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"io/ioutil"

	"github.com/dolmen-go/hexenc"
)

func readAllByChunks(r io.Reader, expected int, chunks []int) error {
	max := 0
	for _, c := range chunks {
		if c > max {
			max = c
		}
	}
	if max == 0 {
		chunks = []int{expected}
		max = expected
	}

	b := make([]byte, max)
	i := 0
	for expected > 0 {
		b = b[:chunks[i]]
		if len(b) > expected {
			b = b[:expected]
		}
		n, err := r.Read(b)
		expected -= n
		if err == io.EOF {
			if expected != 0 {
				return fmt.Errorf("EOF but expecting %d more bytes", expected)
			}
			// Try again: this should give the same error
			n, err := r.Read(b)
			if err != io.EOF {
				return fmt.Errorf("Got EOF once, but not twice: %v", err)
			}
			if n != 0 {
				return fmt.Errorf("After EOF n should be 0 but is %d", n)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Unexpected read error: %s", err)
		}
		if n == 0 && len(b) > 0 {
			return errors.New("Read error: got 0 bytes")
		}
		i = (i + 1) % len(chunks)
	}
	return nil
}

type ChunkedReader struct {
	r       io.Reader
	current int
	sizes   []int
}

func (r *ChunkedReader) Read(b []byte) (n int, err error) {
	size := r.sizes[r.current]
	r.current = (r.current + 1) % len(r.sizes)
	if len(b) > size {
		b = b[:size]
	}
	return r.r.Read(b)
}

func NewChunkedReader(r io.Reader, sizes []int) *ChunkedReader {
	return &ChunkedReader{r: r, sizes: sizes}
}

func TestDecoder(t *testing.T) {
	for _, test := range []struct {
		size  int
		sizes []int
	}{
		{0, nil},
		{1, nil},
		{1023, []int{1}},
		{1023, []int{2}},
		{1023, []int{4}},
		{1023, []int{64}},
		{32768, []int{64}},
		{32768, []int{13}},
		{1024, []int{16}},
		{1024, []int{16, 32, 5}},
		{500, []int{16, 32, 5, 0}},
	} {
		// Make random data
		data := make([]byte, test.size)
		rand.Read(data)
		t.Logf("Test data: %x", data)

		// Encode to hex
		buf := &bytes.Buffer{}
		enc := hexenc.Encoding{}.NewEncoder(buf)
		n, err := enc.Write(data)
		if n != len(data) {
			t.Errorf("Encoding error: %d written", n)
			continue
		}
		if err != nil {
			t.Errorf("Encoding error: %s", err)
			continue
		}

		// Verify encoding against Sprintf %x
		if buf.String() != fmt.Sprintf("%x", data) {
			t.Errorf("Encoding error: no match")
			continue
		}

		hexData := buf.Bytes()

		dec := hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
		b, err := ioutil.ReadAll(dec)
		if err != nil {
			t.Errorf("Read error: %s", err)
			continue
		}
		if !bytes.Equal(b, data) {
			t.Errorf("Decode failure: %x != %x", b, data)
			continue
		}

		if test.sizes != nil {
			// Feed decoder with incomplete chunks
			dec = hexenc.Encoding{}.NewDecoder(NewChunkedReader(bytes.NewReader(hexData), test.sizes))
			b, err = ioutil.ReadAll(dec)
			if err != nil {
				t.Errorf("Read error: %s", err)
				continue
			}
			if !bytes.Equal(b, data) {
				t.Errorf("Decode failure: %x != %x", b, data)
				continue
			}
		}

		dec = hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
		buf = &bytes.Buffer{}
		t.Logf("Read by chunks %v", test.sizes)
		err = readAllByChunks(io.TeeReader(dec, buf), len(data), test.sizes)
		if err != nil {
			t.Errorf("Read error: %s", err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), data) {
			t.Errorf("Decode failure: %x != %x", buf.Bytes(), data)
			continue
		}
	}
}
