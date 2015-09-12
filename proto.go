package proto

import (
	"fmt"
	"io"
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
