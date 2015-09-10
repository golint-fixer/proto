package qik

import (
	"bytes"
	"crypto/rand"
	"github.com/johnmcconnell/nop"
	"github.com/johnmcconnell/proto"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func randomBytes(S int) ([]byte, error) {
	BS := make([]byte, S)

	_, err := rand.Read(BS)

	if err != nil {
		return nil, err
	}

	return BS, err
}

func TestEncoding(t *testing.T) {
	assert := assert.New(t)

	BS := []byte{0, 1, 2, 3, 4}
	expected := []byte{0, 5, 0, 1, 2, 3, 4}

	W := bytes.NewBuffer(nil)
	E := NewWriter(W)

	n, err := E.Write(BS)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		5,
		n,
		"5 bytes were written",
	)

	assert.Equal(
		expected,
		W.Bytes(),
		"bytes match",
	)
}

func TestDecoding(t *testing.T) {
	assert := assert.New(t)

	BS := []byte{0, 5, 0, 1, 2, 3, 4}
	expected := []byte{0, 1, 2, 3, 4}

	var complete []byte

	R := bytes.NewBuffer(BS)
	D := NewReader(R)
	B := make([]byte, 3)

	n, err := D.Read(B)
	complete = append(complete, B...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		3,
		n,
		"3 bytes were read",
	)

	assert.Equal(
		[]byte{0, 1, 2},
		B,
		"bytes match",
	)

	n, err = D.Read(B)
	complete = append(complete, B[:n]...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		2,
		n,
		"2 bytes were read",
	)

	assert.Equal(
		[]byte{3, 4},
		B[:n],
		"bytes match",
	)

	n, err = D.Read(B)

	assert.Equal(
		io.EOF,
		err,
		"Using the proto EOF to designate end of file",
	)

	assert.Equal(
		0,
		n,
		"No bytes were read",
	)

	assert.Equal(
		expected,
		complete,
		"bytes match",
	)
}

func TestDecodingMult(t *testing.T) {
	assert := assert.New(t)

	BS := []byte{0, 5, 0, 1, 2, 3, 4, 0, 3, 0, 1, 2}
	expected := []byte{0, 1, 2, 3, 4, 0, 1, 2}

	var complete []byte

	R := bytes.NewBuffer(BS)
	D := NewReader(R)
	B := make([]byte, 3)

	n, err := D.Read(B)
	complete = append(complete, B...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		3,
		n,
		"3 bytes were read",
	)

	assert.Equal(
		[]byte{0, 1, 2},
		B,
		"bytes match",
	)

	n, err = D.Read(B)
	complete = append(complete, B[:n]...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		2,
		n,
		"2 bytes were read",
	)

	assert.Equal(
		[]byte{3, 4},
		B[:n],
		"bytes match",
	)

	n, err = D.Read(B)
	complete = append(complete, B[:n]...)

	assert.Equal(
		nil,
		err,
		"There is no error",
	)

	assert.Equal(
		3,
		n,
		"3 bytes were read",
	)

	assert.Equal(
		[]byte{0, 1, 2},
		B[:n],
		"bytes match",
	)

	n, err = D.Read(B)

	assert.Equal(
		io.EOF,
		err,
		"Using io.EOFM to designate end of file",
	)

	assert.Equal(
		0,
		n,
		"No bytes were read",
	)

	assert.Equal(
		expected,
		complete,
		"bytes match",
	)
}

func TestDecodingMultMessage(t *testing.T) {
	assert := assert.New(t)

	BS := []byte{0, 2, 0, 1, 0, 0, 0, 2, 0, 1}
	expected := []byte{0, 1, 0, 1}

	var complete []byte

	R := bytes.NewBuffer(BS)
	D := NewReader(R)
	B := make([]byte, 3)

	n, err := D.Read(B)
	complete = append(complete, B[:n]...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		2,
		n,
		"2 bytes were read",
	)

	assert.Equal(
		[]byte{0, 1},
		B[:n],
		"bytes match",
	)

	n, err = D.Read(B)

	assert.Equal(
		proto.ErrEOM,
		err,
		"Using proto to designate end of message",
	)

	n, err = D.Read(B)
	complete = append(complete, B[:n]...)

	assert.Nil(
		err,
		"Error is nil",
	)

	assert.Equal(
		2,
		n,
		"2 bytes were read",
	)

	assert.Equal(
		[]byte{0, 1},
		B[:n],
		"bytes match",
	)

	n, err = D.Read(B)

	assert.Equal(
		io.EOF,
		err,
		"Using io.EOF to designate end of file",
	)

	assert.Equal(
		0,
		n,
		"No bytes were read",
	)

	assert.Equal(
		expected,
		complete,
		"bytes match",
	)
}

func encodeBenchmarkSerial(size int, b *testing.B) {
	bs, _ := randomBytes(size)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		W := nop.NewWriter()
		E := NewWriter(W)

		b.StartTimer()
		E.Write(bs)
	}
}

func BenchmarkEncode_100B(b *testing.B) {
	encodeBenchmarkSerial(100, b)
}

func BenchmarkEncode_1K(b *testing.B) {
	encodeBenchmarkSerial(1000, b)
}

func BenchmarkEncode_100K(b *testing.B) {
	encodeBenchmarkSerial(100*1000, b)
}

func BenchmarkEncode_1M(b *testing.B) {
	encodeBenchmarkSerial(1*1000*1000, b)
}

func BenchmarkEncode_100M(b *testing.B) {
	encodeBenchmarkSerial(100*1000*1000, b)
}

func decodeBenchmarkSerial(size int, b *testing.B) {
	bs, _ := randomBytes(size)
	buff := make([]byte, 512)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		R := bytes.NewBuffer(bs)
		D := NewReader(R)

		b.StartTimer()
		D.Read(buff)
	}
}

func BenchmarkDecode_100B(b *testing.B) {
	decodeBenchmarkSerial(100, b)
}

func BenchmarkDecode_1K(b *testing.B) {
	decodeBenchmarkSerial(1000, b)
}

func BenchmarkDecode_100K(b *testing.B) {
	decodeBenchmarkSerial(100*1000, b)
}

func BenchmarkDecode_1M(b *testing.B) {
	decodeBenchmarkSerial(1*1000*1000, b)
}

func BenchmarkDecode_100M(b *testing.B) {
	decodeBenchmarkSerial(100*1000*1000, b)
}
