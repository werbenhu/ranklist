<div align='center'>
<a href="https://github.com/werbenhu/ranklist/actions"><img src="https://github.com/werbenhu/ranklist/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/ranklist"><img src="https://goreportcard.com/badge/github.com/werbenhu/ranklist"></a>
<a href="https://coveralls.io/github/werbenhu/ranklist?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/ranklist/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/ranklist"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/ranklist"><img src="https://pkg.go.dev/badge/github.com/werbenhu/ranklist.svg"></a>
</div>

[English](README.md) | [简体中文](README_CN.md)

# ranklist

A high-performance real-time ranking data structure implemented using a Skip List in Golang.

## Features

- Thread-Safe Operations: Provides safe concurrent access.
- Generic Support: Works seamlessly with various comparable data types.
- O(log n) Time Complexity: Efficient insertion, deletion, and query operations.
- Real-Time Ranking Queries: Optimized for fast ranking updates.
- Secondary Sorting: Supports tie-breaking for equal scores.
- Built-In Key-Value Dictionary: Enables O(1) key-value lookups.

### Exceptional Performance

Whether for real-time ranking systems or as a high-performance key-value storage, ranklist delivers outstanding efficiency, achieving millions of writes and reads per second effortlessly.

- Write Performance: Capable of handling over a million write requests per second.
- Read Performance: Can handle up to tens of millions of read operations per second.
- Ranking Query: Supports real-time ranking queries, with the ability to process millions of ranking queries per second.

```
goos: windows
goarch: amd64
pkg: github.com/werbenhu/ranklist
cpu: AMD Ryzen 5 5600H with Radeon Graphics
BenchmarkRankListSet-12                  2656849               413.0 ns/op           288 B/op          1 allocs/op
BenchmarkRankListRandSet-12              1000000              1679 ns/op             309 B/op          1 allocs/op
BenchmarkRankListGet-12                 14207959                87.78 ns/op            0 B/op          0 allocs/op
BenchmarkRankListRank-12                 6200530               192.3 ns/op             0 B/op          0 allocs/op
BenchmarkRankListRange-12                5907865               203.3 ns/op           496 B/op          5 allocs/op
BenchmarkZSetRandSet-12                  1000000              3302 ns/op             167 B/op          3 allocs/op
BenchmarkFastSkipListSet-12              6264999               207.2 ns/op            68 B/op          2 allocs/op
BenchmarkFastSkipListRandSet-12          1000000              1556 ns/op              68 B/op          2 allocs/op
BenchmarkFastSkipListGet-12             13096218                91.68 ns/op            0 B/op          0 allocs/op
BenchmarkMapSet-12                       4131572               365.1 ns/op           127 B/op          1 allocs/op
PASS
ok      github.com/werbenhu/ranklist    25.241s
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/werbenhu/ranklist"
)

func main() {
	// Create a new rank list where keys are strings and scores are integers.
	// The rank list uses a skip list internally for efficient ranking operations.
	r := ranklist.New[string, int]()

	// Add elements to the rank list with their respective scores.
	// Keys "a", "b", "c", "d", and "e" are assigned scores 1, 2, 3, 4, and 5, respectively.
	r.Set("a", 1)
	r.Set("b", 2)
	r.Set("c", 3)
	r.Set("d", 4)
	r.Set("e", 5)

	// Delete the key "e" from the rank list.
	// The Del method returns true if the key existed and was successfully removed.
	if ok := r.Del("e"); ok {
		fmt.Printf("Successfully deleted 'e'\n")
	}

	// Get the rank of the key "c".
	// The Rank method returns the rank of the key (1-based) and a boolean indicating success.
	if rank, ok := r.Rank("c"); ok {
		fmt.Printf("The rank of 'c' is: %d\n", rank)
	}

	// Get the score associated with the key "d".
	// The Get method returns the score and a boolean indicating success.
	if score, ok := r.Get("d"); ok {
		fmt.Printf("The score of 'd' is: %d\n", score)
	}

	// Retrieve the top 3 keys and their scores from the rank list.
	ranks := r.Range(1, 4)
	startRank := 1
	for k := range ranks {
		fmt.Printf("Key: %s, Score: %d, Rank: %d\n", ranks[k].Key, ranks[k].Value, startRank)
		startRank++
	}
}
```