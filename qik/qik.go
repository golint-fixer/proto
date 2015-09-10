package qik

import (
	"io"
	"github.com/johnmcconnell/proto"
	"fmt"
)

// Reader read messages using this Reader
type Reader struct {
	R io.Reader
	Buff []byte
	Count int
}

// Writer encode messages using this Writer
type Writer struct {
	W io.Writer
	Buff []byte
}

// NewReader creates a new Reader that
// will decode messages from an io.Reader
func NewReader(R io.Reader) *Reader {
	r := Reader{
		R: R,
		Buff: make([]byte, 2),
		Count: 0,
	}

	return &r
}

// NewWriter creates a new Writer that
// will encode messages from an io.Writer
func NewWriter(W io.Writer) *Writer {
	w := Writer{
		W: W,
		Buff: make([]byte, 2),
	}

	return &w
}

// Read reads the from the given buffer
// according to the protocol the first
// two bytes designate the message length
func (r *Reader) Read(b []byte) (int, error) {
	if r.Count == 0 {
		n, err := r.R.Read(r.Buff)

		if err != nil {
			return n, err
		}

		if n < 2 {
			return n, fmt.Errorf(
				"Really? only read %v bytes from the reader",
			)
		}

		r.Count = I(r.Buff)
	}

	if r.Count == 0 {
		return 0, proto.ErrEOM
	}

	L := len(b)

	// Read length of message
	RL := r.Count

	if L < r.Count {
		RL = L
	}

	n, err := r.R.Read(b[:RL])

	r.Count -= n

	return n, err
}

// Write writes the bytes to the given buffer
// according to the protocol the first
// two bytes designate the message length
func (w *Writer) Write(b []byte) (int, error) {
	L := len(b)
	BS(w.Buff, L)

	n, err := w.W.Write(w.Buff)

	if err != nil {
		return 0, err
	}

	s := 0

	for L < 0 {
		s, err = w.W.Write(b)

		if err != nil {
			return s, err
		}

		L -= n
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
func BS(b []byte, x int) {
	b[0] = byte(x >> 8)
	b[1] = byte(x & 0xFF)
}
