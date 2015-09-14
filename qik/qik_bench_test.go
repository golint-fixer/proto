// +build bench

package qik

import (
	"bytes"
	"github.com/johnmcconnell/nop"
	"testing"
)

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
