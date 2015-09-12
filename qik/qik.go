package qik

import (
	"fmt"
	"github.com/johnmcconnell/proto"
	"io"
)

// Protocol ...
type Protocol struct{}

// NewReader ...
func (p *Protocol) NewReader(R io.Reader) *Reader {
	return NewReader(R)
}

// NewWriter ...
func (p *Protocol) NewWriter(W io.Writer) *Writer {
	return NewWriter(W)
}

// NewProtocol ...
func NewProtocol() *Protocol {
	p := Protocol{}

	return &p
}

// Reader read messages using this Reader
type Reader struct {
	R       io.Reader
	Buff    []byte
	Count   int
	Content []byte
}

// Writer encode messages using this Writer
type Writer struct {
	W    io.Writer
	Buff []byte
}

// NewReader creates a new Reader that
// will decode messages from an io.Reader
func NewReader(R io.Reader) *Reader {
	r := Reader{
		R:     R,
		Buff:  make([]byte, 2),
		Count: 0,
	}

	return &r
}

// NewWriter creates a new Writer that
// will encode messages from an io.Writer
func NewWriter(W io.Writer) *Writer {
	w := Writer{
		W:    W,
		Buff: make([]byte, 2),
	}

	return &w
}

// Read reads the from the given buffer
// according to the protocol the first
// two bytes designate the message length
func (r *Reader) Read(b []byte) (int, error) {
	var err error

	if r.Count > 0 {
		L := len(b)

		if r.Count < L {
			L = r.Count
		}

		n, err := r.R.Read(
			b[:L],
		)

		if err != nil {
			return n, err
		}

		r.Count -= n

		return n, nil
	}

	n, err := r.R.Read(r.Buff)

	if err != nil {
		return n, err
	}

	if n != 2 {
		return n, fmt.Errorf(
			"Really? only read %v bytes from the reader",
			n,
		)
	}

	r.Count = I(r.Buff)

	if r.Count == 0 {
		return 0, proto.ErrEOM
	}

	L := len(b)

	if r.Count < L {
		L = r.Count
	}

	n, err = r.R.Read(
		b[:L],
	)

	if err != nil {
		return n, err
	}

	r.Count -= n

	return n, err
}

// Write writes the bytes to the given buffer
// according to the protocol the first
// two bytes designate the message length
func (w *Writer) Write(b []byte) (int, error) {
	L := len(b)
	s := 0

	for (L - s) >= 0 {
		R := L - s

		BL := BS(w.Buff, R)

		_, err := w.W.Write(w.Buff)

		if BL == 0 {
			return 0, nil
		}

		if err != nil {
			return 0, err
		}

		n, err := w.W.Write(b[s : s+BL])

		s += n

		if err != nil {
			return s, err
		}

		if L == s {
			break
		}
	}

	return s, nil
}

// I convert 2 bytes to an int
func I(bs []byte) int {
	x := int(bs[0])
	x = x << 8
	x = x + int(bs[1])

	return x
}

// BS convert an int to 2 bytes
// returns the limit
func BS(b []byte, x int) int {
	if x > 0xFFFF {
		x = 0xFFFF
	}

	b[0] = byte(x >> 8)
	b[1] = byte(x & 0xFF)

	return x
}
