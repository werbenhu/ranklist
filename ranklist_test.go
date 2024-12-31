package ranklist

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	sl := New[int, int]()

	if sl.length != 0 {
		t.Errorf("New skiplist should have length 0, got %d", sl.length)
	}
	if sl.level != 1 {
		t.Errorf("New skiplist should have initial level 1, got %d", sl.level)
	}
}

func TestSetHighProbability(t *testing.T) {
	PROBABILITY = 0.7
	sl := New[string, int]()
	for k := 0; k < 10000; k++ {
		sl.Set(strconv.Itoa(k), k)
	}
}

func TestSetLowProbability(t *testing.T) {
	PROBABILITY = 0.05
	sl := New[string, int]()
	for k := 0; k < 10000; k++ {
		sl.Set(strconv.Itoa(k), k)
	}
}

func TestSetAndGet(t *testing.T) {
	sl := New[string, int]()

	sl.Set("a", 1)
	sl.Set("b", 2)
	sl.Set("c", 3)

	testCases := []struct {
		key      string
		expected int
		exists   bool
	}{
		{"a", 1, true},
		{"b", 2, true},
		{"c", 3, true},
		{"d", 0, false},
	}

	for _, tc := range testCases {
		value, exists := sl.Get(tc.key)
		if exists != tc.exists {
			t.Errorf("Key %s: unexpected exists status, got %v", tc.key, exists)
		}
		if exists && value != tc.expected {
			t.Errorf("Key %s: expected value %d, got %d", tc.key, tc.expected, value)
		}
	}
}

func TestUpdate(t *testing.T) {
	sl := New[int, int]()

	sl.Set(1, 100)
	sl.Set(1, 200)

	value, exists := sl.Get(1)
	if !exists {
		t.Fatalf("Key should exist after update")
	}
	if value != 200 {
		t.Errorf("Value should be updated to 200, got %d", value)
	}
}

func TestDel(t *testing.T) {
	sl := New[int, int]()

	sl.Set(1, 100)
	sl.Set(2, 200)
	sl.Set(2, 300)

	if !sl.Del(1) {
		t.Fatalf("Delete should return true for existing key")
	}
	if _, exists := sl.Get(1); exists {
		t.Fatalf("Key should not exist after deletion")
	}
	if sl.Del(3) {
		t.Fatalf("Delete should return false for non-existent key")
	}
}

func TestRank(t *testing.T) {
	sl := New[string, int]()

	testData := []struct {
		key   string
		value int
		rank  int
	}{
		{"a", 1, 1},
		{"b", 2, 2},
		{"c", 3, 3},
		{"d", 3, 4},
		{"e", 4, 5},
	}

	for _, data := range testData {
		sl.Set(data.key, data.value)
	}

	for _, data := range testData {
		rank, exists := sl.Rank(data.key)
		if !exists {
			t.Fatalf("Key %s should exist", data.key)
		}
		if rank != data.rank {
			t.Errorf("Key %s: expected rank %d, got %d", data.key, data.rank, rank)
		}
	}

	rank, exists := sl.Rank("x")
	if exists {
		t.Fatalf("Rank should return false for non-existent key")
	}
	if rank != 0 {
		t.Errorf("Rank should be 0 for non-existent key, got %d", rank)
	}
}

func TestMassiveRank(t *testing.T) {
	sl := New[string, int64]()
	for k := 0; k < 10000; k++ {
		sl.Set(strconv.Itoa(k), rand.Int64N(10000))
	}

	for k := 0; k < 10000; k++ {
		_, exists := sl.Rank(strconv.Itoa(k))
		if !exists {
			t.Errorf("Key %s should exist", strconv.Itoa(k))
		}
	}
}

func TestRankKeyNotExist(t *testing.T) {
	sl := New[string, int]()
	sl.dict["x"] = nil

	rank, exists := sl.Rank("x")
	if exists {
		t.Fatalf("Rank should return false for non-existent key")
	}
	if rank != 0 {
		t.Errorf("Rank should be 0 for non-existent key, got %d", rank)
	}
}

func TestEdgeCases(t *testing.T) {
	sl := New[int, int]()

	if _, exists := sl.Get(1); exists {
		t.Fatalf("Get should return false for empty skip list")
	}
	if sl.Del(1) {
		t.Fatalf("Delete should return false for empty skip list")
	}

	for i := 0; i < 1000; i++ {
		sl.Set(i, i)
	}

	if sl.length != 1000 {
		t.Errorf("Length should be 1000 after insertions, got %d", sl.length)
	}

	for i := 0; i < 1000; i++ {
		if !sl.Del(i) {
			t.Errorf("Failed to delete key %d", i)
		}
	}

	if sl.length != 0 {
		t.Errorf("Length should be 0 after all deletions, got %d", sl.length)
	}
}
