# proto [![Build Status](https://travis-ci.org/johnmcconnell/proto.svg?branch=master)](https://travis-ci.org/johnmcconnell/proto)
Performant Message Protocols over TCP keep alive

## Performance

### Qik Protocol

```
BenchmarkEncode_100B  10000000         204 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1K  10000000         204 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100K  10000000         183 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1M   3000000         397 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100M     50000       24926 ns/op        49 B/op        1 allocs/op
BenchmarkDecode_100B   3000000         526 ns/op       194 B/op        2 allocs/op
BenchmarkDecode_1K   3000000         548 ns/op       194 B/op        2 allocs/op
BenchmarkDecode_100K   5000000         366 ns/op       194 B/op        2 allocs/op
BenchmarkDecode_1M   5000000         262 ns/op       193 B/op        2 allocs/op
BenchmarkDecode_100M   5000000         243 ns/op       193 B/op        2 allocs/op
ok    github.com/johnmcconnell/proto/qik  91.255s
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
