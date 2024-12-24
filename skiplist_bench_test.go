package skiplist

import (
	"math/rand/v2"
	"strconv"
	"testing"

	fastskiplist "github.com/sean-public/fast-skiplist"
)

// Insert keys randomly
func BenchmarkSkipListInsertRandom(b *testing.B) {
	list := NewSkiplist[float64, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Set(rand.Float64(), i)
	}
}

// Insert keys in sorted order
func BenchmarkSkipListInsertSorted(b *testing.B) {
	list := NewSkiplist[float64, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Set(float64(i), i)
	}
}

// Insert keys randomly
func BenchmarkFastSkipListInsertRandom(b *testing.B) {
	fast := fastskiplist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Set(rand.Float64(), i)
	}
}

// Insert keys in sorted order
func BenchmarkFastSkipListInsertSorted(b *testing.B) {
	fast := fastskiplist.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fast.Set(float64(i), i)
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
