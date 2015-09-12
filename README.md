# proto [![Build Status](https://travis-ci.org/johnmcconnell/proto.svg?branch=master)](https://travis-ci.org/johnmcconnell/proto)
Performant Message Protocols over TCP keep alive

## Performance

### [Qik Protocol](qik)

```
PASS
BenchmarkEncode_100B  10000000         198 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1K  10000000         201 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100K  10000000         168 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_1M   3000000         482 ns/op        50 B/op        1 allocs/op
BenchmarkEncode_100M    200000       10169 ns/op        49 B/op        1 allocs/op
BenchmarkDecode_100B   3000000         473 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_1K   3000000         478 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_100K   5000000         335 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_1M   5000000         244 ns/op       162 B/op        2 allocs/op
BenchmarkDecode_100M   5000000         226 ns/op       161 B/op        2 allocs/op
ok    github.com/johnmcconnell/proto/qik  94.951s
```

### [Slim Protocol](slim)
