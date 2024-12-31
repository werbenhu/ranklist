<div align='center'>
<a href="https://github.com/werbenhu/ranklist/actions"><img src="https://github.com/werbenhu/ranklist/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/ranklist"><img src="https://goreportcard.com/badge/github.com/werbenhu/ranklist"></a>
<a href="https://coveralls.io/github/werbenhu/ranklist?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/ranklist/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/ranklist"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/ranklist"><img src="https://pkg.go.dev/badge/github.com/werbenhu/ranklist.svg"></a>
</div>

[English](README.md) | [简体中文](README_CN.md)

# ranklist

A high-performance ranking system implementation based on Skip List.

## Features

- Thread-Safe Operations: Provides safe concurrent access.
- Generic Support: Works seamlessly with various comparable data types.
- O(log n) Time Complexity: Efficient insertion, deletion, and query operations.
- Real-Time Ranking Queries: Optimized for fast ranking updates.
- Secondary Sorting: Supports tie-breaking for equal scores.
- Built-In Key-Value Dictionary: Enables O(1) key-value lookups.

### Exceptional Performance

Whether for real-time ranking systems or as a high-performance key-value storage, ranklist delivers outstanding efficiency, achieving millions of writes and reads per second effortlessly.

- Write Performance: Each write operation takes just 812.4 nanoseconds, handling over 1 million writes per second.
- Read Performance: Each read operation takes only 64.03 nanoseconds, achieving 10+ million reads per second.
- Ranking Queries: Real-time ranking queries execute in 377.4 nanoseconds, completing over 2 million calculations per second.

```
goos: windows
goarch: amd64
pkg: github.com/werbenhu/ranklist
cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
BenchmarkRankListSet-16          1451098               812.4 ns/op           445 B/op          5 allocs/op
BenchmarkRankListGet-16         18473389                64.03 ns/op            0 B/op          0 allocs/op
BenchmarkRankListRank-16         3613449               377.4 ns/op             0 B/op          0 allocs/op
BenchmarkFastSkipListSet-16      1000000              1116 ns/op              68 B/op          2 allocs/op
BenchmarkMapSet-16               5044543               271.7 ns/op           106 B/op          1 allocs/op
PASS
ok      github.com/werbenhu/ranklist    17.302s
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/werbenhu/ranklist"
)

func main() {
    r := ranklist.New[string, int]()

    r.Set("a", 1)
    r.Set("b", 2)
    r.Set("c", 3)
    r.Set("d", 4)
    r.Set("e", 5)

    if ok := r.Del("e"); ok {
        fmt.Printf("Successfully deleted 'e'\n")
    }

    if rank, ok := r.Rank("c"); ok {
        fmt.Printf("The rank of 'c' is: %d\n", rank)
    }

    if score, ok := r.Get("d"); ok {
        fmt.Printf("The score of 'd' is: %d\n", score)
    }
}
```