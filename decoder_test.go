package hexenc_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/dolmen-go/hexenc"
)

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

		// Try to read the whole stream in one Read call.
		// My encoder can do it, but the one from stdlib (in go1.10) doesn't.
		dec := hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
		b := make([]byte, len(data))
		n, err = dec.Read(b)
		if err != nil && err != io.EOF {
			t.Errorf("Read error: %s", err)
			continue
		}
		if n != len(b) {
			t.Logf("Incomplete read: got %d bytes, expected %d", n, len(b))
		}
		if !bytes.Equal(b[:n], data[:n]) {
			t.Errorf("Decode failure: %x != %x", b[:n], data[:n])
			continue
		}

		checkReadAll := func(decoder io.Reader) error {
			b, err := ioutil.ReadAll(decoder)
			if err != nil {
				return fmt.Errorf("Read error: %s", err)
			}
			if !bytes.Equal(b, data) {
				return fmt.Errorf("Decode failure: %x != %x", b, data)
			}
			return nil
		}

		dec = hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
		if err := checkReadAll(dec); err != nil {
			t.Error(err)
			continue
		}

		if test.sizes != nil {
			// Feed decoder with incomplete chunks
			t.Logf("Feed decoder by chunks %v", test.sizes)
			dec = hexenc.Encoding{}.NewDecoder(NewChunkedReader(bytes.NewReader(hexData), test.sizes))
			if err := checkReadAll(dec); err != nil {
				t.Error(err)
				continue
			}

			t.Logf("Read from decoder by chunks %v", test.sizes)
			dec = hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
			if err := checkReadAll(NewChunkedReader(dec, test.sizes)); err != nil {
				t.Error(err)
				continue
			}
		}
	}
}
