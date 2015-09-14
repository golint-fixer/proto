package proto

import (
	"fmt"
	"io"
	"net"
)

var (
	// ErrEOM end of message error
	ErrEOM = fmt.Errorf(
		"EOM",
	)
)

// Protocol is an interface to build the actual
// underlying message readers and writers. It is
// important to be able to swap out protocols
// with others
type Protocol interface {
	NewReader(io.Reader) io.Reader
	NewWriter(io.Writer) io.Writer
}

// Conn is a net.Conn but
// read and writing bytes goes
// through the protocol
type Conn struct {
	net.Conn
	W io.Writer
	R io.Reader
}

// Write write from the connection but only
// after going through the protocol
func (c *Conn) Write(b []byte) (int, error) {
	return c.W.Write(b)
}

// Read read from the connection but only
// after going through the protocol
func (c *Conn) Read(b []byte) (int, error) {
	return c.R.Read(b)
}

// MessageReader has the exact same interface as
// io.Reader but can also return a proto.EOM error
// which designates the end of a message but not the
// stream
type MessageReader interface {
	Read([]byte) (int, error)
}

// MessageWriter has the exat same interface as
// io.Writer. Writing the nil buffer or the empty
// buffer to a MessageWriter designates the end
// of a message and the start of a new one
type MessageWriter interface {
	Write([]byte) (int, error)
}

// WrapConn wraps a network connection around
// a protocol
func WrapConn(p Protocol, c net.Conn) net.Conn {
	C := Conn{
		Conn: c,
		W:    p.NewWriter(c),
		R:    p.NewReader(c),
	}

	return &C
}

// CopyMessages copy N messages from the reader
// to the writer, returns any errors.
// returns the number of messages written successful
// otherwise it will return a non-nil error and the
// number of bytes in the buffer
func CopyMessages(W io.Writer, R io.Reader, B []byte, N int) (int, error) {
	i := 0
	infinite := N < 0

	for i < N || infinite {
		n, err := R.Read(B)

		if err == ErrEOM {
			i++

			// Don't send an extra end of message markers
			if n != 0 {
				_, err := W.Write(B[:n])

				if err != nil {
					return n, err
				}
			}

			// End of message marker
			_, err := W.Write(nil)

			if err != nil {
				return n, err
			}

			continue
		}

		if err == io.EOF {
			return n, err
		}

		if err != nil {
			return n, err
		}

		_, err = W.Write(B[:n])

		if err != nil {
			return n, err
		}
	}

	return i, nil
}

// WriteMessage ...
func WriteMessage(W io.Writer, Bytes []byte) (int, error) {
	S := 0

	n, err := W.Write(Bytes)

	S += n

	if err != nil {
		return S, err
	}

	_, err = W.Write(nil)

	if err != nil {
		return S, err
	}

	return S, nil
}

// ReadMessage ...
func ReadMessage(R io.Reader) ([]byte, error) {
	var B []byte

	l := 16
	for {
		b := make([]byte, l)

		n, err := R.Read(b)

		if n == l {
			l *= 2
		}

		B = append(B, b[:n]...)

		if err == io.EOF {
			return B, nil
		}

		if err == ErrEOM {
			break
		}

		if err != nil {
			return B, err
		}
	}

	return B, nil
}

// Copy copy the bytes into a if a cannot
// hold b, grow a in order to contain it
func Copy(a, b []byte) []byte {
	La := len(a)
	Lb := len(b)

	if La < Lb {
		a = make([]byte, Lb)
	}

	n := copy(a, b)

	return a[:n]
}
