package ranklist

import (
	"math/rand/v2"
	"strconv"
	"testing"

	fastskiplist "github.com/sean-public/fast-skiplist"
)

func BenchmarkRankListSet(b *testing.B) {
	sl := New[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
}

func BenchmarkRankListGet(b *testing.B) {
	sl := New[int, int]()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(i % 1000000)
	}
}

func BenchmarkRankListRank(b *testing.B) {
	sl := New[int, int]()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Rank(i % 1000000)
	}
}

func BenchmarkRankListRange(b *testing.B) {
	sl := New[int, int]()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Range(0, 10)
	}
}

func BenchmarkFastSkipListSet(b *testing.B) {
	fast := fastskiplist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Set(rand.Float64(), i)
	}
}

func BenchmarkMapSet(b *testing.B) {
	n := make(map[string]float64, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		score := rand.Float64() * 1000
		n[strconv.Itoa(i)] = score
	}
}
