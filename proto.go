package proto

import (
	"fmt"
)

var (
	// ErrEOM end of message error
	ErrEOM = fmt.Errorf(
		"EOM",
	)
)

// MessageReader has the exact same interface as
// io.Reader but can also return a proto.EOM error
// which designates the end of a message but not the
// stream
type MessageReader interface {
	Read([]byte) (int, error)
}

// MessageWriter has the exat same interface as
// io.Writer but also includes EOM. EOM designates
// that everything written is a single message and
// everything written afterward will be a new message
type MessageWriter interface {
	Write([]byte) (int, error)
	EOM()
}

// Copy copy the bytes into a if a cannot
// hold b, grow a in order to contain it
func Copy(a, b []byte) []byte {
	La := len(a)
	Lb := len(b)

	if La < Lb {
		a = make([]byte, Lb)
	}

	copy(a, b)

	return a
}
