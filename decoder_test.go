package hexenc_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"io/ioutil"

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

		dec = hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
		b = make([]byte, len(data))
		n, err = dec.Read(b)
		if err != nil && err != io.EOF {
			t.Errorf("Read error: %s", err)
			continue
		}
		if !bytes.Equal(b, data) {
			t.Errorf("Decode failure: %x != %x", b, data)
			continue
		}

		if test.sizes != nil {
			// Feed decoder with incomplete chunks
			t.Logf("Feed decoder by chunks %v", test.sizes)
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

			t.Logf("Read from decoder by chunks %v", test.sizes)
			dec = hexenc.Encoding{}.NewDecoder(bytes.NewReader(hexData))
			b, err = ioutil.ReadAll(NewChunkedReader(dec, test.sizes))
			if err != nil {
				t.Errorf("Read error: %s", err)
				continue
			}
			if !bytes.Equal(b, data) {
				t.Errorf("Decode failure: %x != %x", buf.Bytes(), data)
				continue
			}
		}
	}
}
