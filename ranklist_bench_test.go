package ranklist

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/liyiheng/zset"
	skiplist "github.com/sean-public/fast-skiplist"
)

func BenchmarkRankListSet(b *testing.B) {
	sl := New[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
}

func BenchmarkRankListRandSet(b *testing.B) {
	sl := New[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, rand.Int())
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

func BenchmarkZSetRandSet(b *testing.B) {
	s := zset.New[int64]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set(rand.Float64(), int64(i))
	}
}

func BenchmarkFastSkipListSet(b *testing.B) {
	fast := skiplist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Set(float64(i), i)
	}
}

func BenchmarkFastSkipListRandSet(b *testing.B) {
	fast := skiplist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Set(rand.Float64(), i)
	}
}

func BenchmarkFastSkipListGet(b *testing.B) {
	fast := skiplist.New()
	for i := 0; i < 1000000; i++ {
		fast.Set(float64(i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Get(float64(i % 1000000))
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
