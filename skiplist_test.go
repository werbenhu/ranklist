package skiplist

import (
	"math/rand/v2"
	"testing"
)

// Insert keys randomly
func BenchmarkSkipListInsertRandom(b *testing.B) {
	list := NewSkiplist()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := rand.IntN(b.N)
		list.Insert(key, float64(i))
	}
}

// Insert keys in sorted order
func BenchmarkSkipListInsertSorted(b *testing.B) {
	list := NewSkiplist()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Insert(i, float64(i))
	}
}
