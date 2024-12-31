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
BenchmarkRankListSet-12                  2691972               407.6 ns/op           288 B/op          1 allocs/op
BenchmarkRankListRandSet-12              1000000              1593 ns/op             309 B/op          1 allocs/op
BenchmarkRankListGet-12                 14354341                83.24 ns/op            0 B/op          0 allocs/op
BenchmarkRankListRank-12                 6383806               191.0 ns/op             0 B/op          0 allocs/op
BenchmarkRankListRange-12                6502486               185.7 ns/op           496 B/op          5 allocs/op
BenchmarkZSetRandSet-12                  1000000              2901 ns/op             167 B/op          3 allocs/op
BenchmarkFastSkipListSet-12              6548966               191.5 ns/op            68 B/op          2 allocs/op
BenchmarkFastSkipListRandSet-12          1000000              1321 ns/op              68 B/op          2 allocs/op
BenchmarkFastSkipListGet-12             13239830                90.82 ns/op            0 B/op          0 allocs/op
BenchmarkMapSet-12                       4489778               328.8 ns/op           118 B/op          1 allocs/op
PASS
ok      github.com/werbenhu/ranklist    23.779s
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