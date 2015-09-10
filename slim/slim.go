package slim

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	// BufferSize the count of bytes the buffer can hold
	BufferSize = 512
	// EscapeByte this byte designates the next character
	// is escaped
	EscapeByte = 0xFE
	// TerminalByte this byte designates the end of the
	// message
	TerminalByte = 0xFF
)

// DecodeString ...
func DecodeString(raw string) (string, error) {
	R := strings.NewReader(raw)
	W := bytes.NewBuffer(nil)

	RBS, err := DecodeMessage(R, W)

	if err != nil && err != io.EOF {
		return "", err
	}

	if len(RBS) != 0 {
		return "", fmt.Errorf(
			"Remaining bytes were returned decoding a string\n\t[%v]\n",
			string(RBS),
		)
	}

	return W.String(), err
}

// DecodeMessage reads bytes off the reader and
// cleans the escape bytes
func DecodeMessage(R io.Reader, W io.Writer) ([]byte, error) {
	BS := make([]byte, BufferSize)

	for {
		n, err := R.Read(BS)

		if err == io.EOF {
			MBS, RBS, err := DecodeBytes(BS[:n])

			if err != nil {
				return nil, err
			}

			W.Write(MBS)

			if RBS == nil {
				return nil, fmt.Errorf(
					"EOF, byte remaining is nil\nm:[%v]",
					MBS,
				)
			}

			return *RBS, err
		}

		if err != nil {
			return nil, err
		}

		MBS, RBS, err := DecodeBytes(BS[:n])

		if err != nil {
			return nil, err
		}

		W.Write(MBS)

		if RBS != nil {
			return *RBS, nil
		}
	}
}

// DecodeBytes removes the two escape bytes
// it returns the decoded bytes in the first return arg
// if it runs into the terminal byte it will return
// the undecoded bytes as the second arg
// if incurs an error it will end early and return
// it as the third arg
func DecodeBytes(BS []byte) ([]byte, *[]byte, error) {
	var message []byte
	L := len(BS)

	escapeNext := false

	for i := range BS {
		b := BS[i]

		if !escapeNext {
			switch b {
			case EscapeByte:
				escapeNext = true

			case TerminalByte:
				// minus one because i at the end is L - 1
				RL := (L - 1) - i

				CBS := make([]byte, RL)
				copy(CBS, BS[i+1:])

				copyBytes := &CBS

				return message, copyBytes, nil

			default:
				message = append(message, b)
			}
		} else { // previous was escape character
			switch b {
			case EscapeByte:
				message = append(message, b)
				escapeNext = false

			case TerminalByte:
				message = append(message, b)
				escapeNext = false

			default:
				return nil, nil, fmt.Errorf(
					"Next character was [%v] but previous was an escape character",
					b,
				)
			}
		}

		i++
	}

	return message, nil, nil
}

// EncodeString ...
func EncodeString(raw string) (string, error) {
	R := strings.NewReader(raw)
	W := bytes.NewBuffer(nil)

	err := EncodeMessage(R, W)

	if err != nil && err != io.EOF {
		return "", err
	}

	return W.String(), err
}

// EncodeMessage encode a message off of a reader
func EncodeMessage(R io.Reader, W io.Writer) error {
	BS := make([]byte, BufferSize)

	for {
		n, err := R.Read(BS)

		if err == io.EOF {
			MBS, _ := EncodeBytes(BS[:n])

			if len(MBS) > 0 {
				W.Write(MBS)
			}

			W.Write(
				[]byte{
					TerminalByte,
				},
			)

			return err
		}

		if err != nil {
			return err
		}

		MBS, _ := EncodeBytes(BS[:n])

		W.Write(MBS)
	}
}

// EncodeBytes adds the two escape bytes
// it returns the encoded bytes in the first return
// it returns the size of the bytes in the second return
func EncodeBytes(BS []byte) ([]byte, int) {
	var message []byte
	L := 0

	for i := range BS {
		b := BS[i]

		switch b {
		case EscapeByte:
			message = append(message, EscapeByte, b)

		case TerminalByte:
			message = append(message, EscapeByte, b)

		default:
			message = append(message, b)
		}

		L++
	}

	return message, L
}
