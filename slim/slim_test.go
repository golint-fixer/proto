package slim

import (
	"bytes"
	"crypto/rand"
	"fmt"
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

	R := bytes.NewReader([]byte{1, 2, 3, 4})
	W := bytes.NewBuffer(nil)

	err := EncodeMessage(R, W)

	assert.Equal(
		err,
		io.EOF,
		"Error is io.EOF",
	)

	BS := []byte{1, 2, 3, 4, TerminalByte}

	assert.Equal(
		BS,
		W.Bytes(),
		"bytes match",
	)
}

func TestEncodeBytes(t *testing.T) {
	assert := assert.New(t)

	BS1 := []byte{1, TerminalByte, EscapeByte, 2, 3, 4}

	BS2, _ := EncodeBytes(BS1)

	BS := []byte{1, EscapeByte, TerminalByte, EscapeByte, EscapeByte, 2, 3, 4}

	assert.Equal(
		BS,
		BS2,
		"bytes match",
	)
}

func TestDecoding(t *testing.T) {
	assert := assert.New(t)

	R := bytes.NewReader([]byte{1, 2, 3, 4, TerminalByte})
	W := bytes.NewBuffer(nil)

	RBS, err := DecodeMessage(R, W)

	assert.Equal(
		err,
		nil,
		"Error is nil",
	)

	assert.Len(
		RBS,
		0,
		"There are no remaining bytes",
	)

	BS := []byte{1, 2, 3, 4}

	assert.Equal(
		W.Bytes(),
		BS,
		"bytes match",
	)
}

func TestDecodeBytes(t *testing.T) {
	assert := assert.New(t)

	BS1 := []byte{1, EscapeByte, EscapeByte, 2, 3, EscapeByte, TerminalByte, 4}

	BS2, RBS, err := DecodeBytes(BS1)

	assert.Nil(
		RBS,
		"RBS is nil",
	)

	assert.Equal(
		err,
		nil,
		"Error is nil",
	)

	BS := []byte{1, EscapeByte, 2, 3, TerminalByte, 4}
	assert.Equal(
		BS,
		BS2,
		"bytes match",
	)

	BS1 = []byte{1, EscapeByte, EscapeByte, 2, 3, TerminalByte}

	BS2, RBS, err = DecodeBytes(BS1)

	assert.NotNil(
		RBS,
		"RBS is not nil",
	)

	assert.Len(
		*RBS,
		0,
		"len(RBS) is not zero",
	)

	assert.Equal(
		err,
		nil,
		"Error is nil",
	)

	BS = []byte{1, EscapeByte, 2, 3}

	assert.Equal(
		BS,
		BS2,
		"bytes match",
	)
}

func TestEncodingAndDecoding(t *testing.T) {
	assert := assert.New(t)

	str1 := "hello world I wish to be encoded and decoded."

	str2, err := EncodeString(str1)
	bytes2 := append([]byte(str1), TerminalByte)

	assert.Equal(
		bytes2,
		[]byte(str2),
		"bytes match",
	)

	assert.Equal(
		io.EOF,
		err,
		"Error is EOF",
	)

	str2, err = DecodeString(str2)

	assert.Equal(
		nil,
		err,
		"Error is nil",
	)

	assert.Equal(
		str1,
		str2,
		"The strings match",
	)
}

func TestLargeBytes(t *testing.T) {
	assert := assert.New(t)

	bs, err := randomBytes(1)

	R := bytes.NewReader(bs)
	B1 := bytes.NewBuffer(nil)

	err = EncodeMessage(R, B1)
	assert.Equal(
		io.EOF,
		err,
		"There is no error",
	)

	B2 := bytes.NewBuffer(nil)
	rbs, err := DecodeMessage(B1, B2)

	assert.Len(
		rbs,
		0,
		fmt.Sprintf(
			"There are no remaing bytes\n\t[%v]\n",
			rbs,
		),
	)
	assert.Nil(
		err,
		"There is no error",
	)

	assert.Equal(
		bs,
		B2.Bytes(),
		"The bytes match",
	)
}

func encodeBenchmarkSerial(size int, b *testing.B) {
	bs, _ := randomBytes(100)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		R := bytes.NewReader(bs)
		B := bytes.NewBuffer(nil)

		b.StartTimer()
		EncodeMessage(R, B)
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
	bs, _ := randomBytes(100)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		R := bytes.NewReader(bs)
		B := bytes.NewBuffer(nil)
		EncodeMessage(R, B)

		O := bytes.NewBuffer(nil)
		b.StartTimer()
		DecodeMessage(B, O)
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
