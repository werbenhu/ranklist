<div align='center'>
<a href="https://github.com/werbenhu/ranklist/actions"><img src="https://github.com/werbenhu/ranklist/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/ranklist"><img src="https://goreportcard.com/badge/github.com/werbenhu/ranklist"></a>
<a href="https://coveralls.io/github/werbenhu/ranklist?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/ranklist/badge.svg?branch=master"></a>   
<a href="https://github.com/werbenhu/ranklist"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/ranklist"><img src="https://pkg.go.dev/badge/github.com/werbenhu/ranklist.svg"></a>
</div>

[English](README.md) | [简体中文](README_CN.md)

# ranklist

一个基于跳表(Skip List)实现的高性能排行榜数据结构。

## 特性 | Features

- 线程安全的操作接口
- 支持泛型，可用于各种可比较的数据类型
- O(log n) 的时间复杂度用于插入、删除和查询操作
- 支持快速的实时排名查询
- 支持相同分数下的二级排序
- 内置键值对字典，提供 O(1) 的键值查找

### 极致性能表现

无论是用于实时排名系统，还是单纯作为高效的 key-value 键值对存储，ranklist 都表现出卓越的性能，轻松实现每秒百万级别的写入与读取操作。

- 写入性能：单次写入仅需 812.4 纳秒，每秒可处理超过百万次写入请求。
- 读取性能：单次读取仅需 64.03 纳秒，每秒读取次数高达千万次以上。
- 排名查询：支持实时排名查询，每次查询仅需 377.4 纳秒，每秒可完成200多万次排名计算。

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

## 使用示例 | Usage

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