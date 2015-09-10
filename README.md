# proto [![Build Status](https://travis-ci.org/johnmcconnell/proto.svg?branch=master)](https://travis-ci.org/johnmcconnell/proto)
Performant Message Protocols over TCP keep alive

## Performance

### Qik Protocol

```
PASS
BenchmarkEncode_100B  10000000         190 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1K  10000000         190 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100K  10000000         155 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1M  10000000         127 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100M  10000000         123 ns/op        49 B/op        1 allocs/op
BenchmarkDecode_100B   2000000         860 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_1K   1000000        1059 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_100K   3000000         600 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_1M   5000000         307 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_100M   5000000         246 ns/op       161 B/op        2 allocs/op
ok    github.com/johnmcconnell/proto/qik  97.635s
```

## Usage

```go
import (
  "github.com/johnmcconnell/proto/qik"
  "bytes"
)


buffer := bytes.NewBuffer(nil)
encoder := qik.NewWriter(buffer)

encoder.Write([]byte("Hello World!")) // first message "Hello World!"
encoder.Write(nil)                    // end of first message
encoder.Write([]byte("Hello World!")) // second message "Hello World!"

bytes := make([]byte, 128)
decoder := qik.NewReader(buffer)

n, _ := decoder.Read(bytes)
string(bytes[:n])                     //=> "Hello World!"

_, err := decoder.Read(bytes)
err                                   //=> proto.ErrEOM{err: "EOM"}

n, _ := decoder.Read(bytes)
string(bytes[:n])                     //=> "Hello World!"

_, err := decoder.Read(bytes)
err                                   //=> io.EOF{err: "EOF"}
```
