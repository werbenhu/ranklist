package ranklist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 测试创建新的跳表
// Test creating new skip list
func TestNew(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()

	r.Equal(0, sl.length, "New skiplist should have length 0")
	r.Equal(1, sl.level, "New skiplist should have initial level 1")
}

// 测试基本的插入和查询操作
// Test basic insertion and query operations
func TestSetAndGet(t *testing.T) {
	r := require.New(t)
	sl := New[string, int]()

	// 测试插入操作
	// Test insertion
	sl.Set("a", 1)
	sl.Set("b", 2)
	sl.Set("c", 3)

	// 测试查询操作
	// Test query
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
		r.Equal(tc.exists, exists, "Key %s: unexpected exists status", tc.key)
		if tc.exists {
			r.Equal(tc.expected, value, "Key %s: unexpected value", tc.key)
		}
	}
}

// 测试更新操作
// Test update operations
func TestUpdate(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()

	// 插入并更新值
	// Insert and update value
	sl.Set(1, 100)
	sl.Set(1, 200)

	value, exists := sl.Get(1)
	r.True(exists, "Key should exist after update")
	r.Equal(200, value, "Value should be updated to 200")
}

// 测试删除操作
// Test deletion operations
func TestDel(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()

	// 插入测试数据
	// Insert test data
	sl.Set(1, 100)
	sl.Set(2, 200)

	// 测试删除存在的键
	// Test deleting existing key
	r.True(sl.Del(1), "Delete should return true for existing key")

	// 验证删除后的状态
	// Verify state after deletion
	_, exists := sl.Get(1)
	r.False(exists, "Key should not exist after deletion")

	// 测试删除不存在的键
	// Test deleting non-existent key
	r.False(sl.Del(3), "Delete should return false for non-existent key")
}

// 测试排名操作
// Test ranking operations
func TestRank(t *testing.T) {
	r := require.New(t)
	sl := New[string, int]()

	// 插入有序数据
	// Insert ordered data
	testData := []struct {
		key   string
		value int
		rank  int
	}{
		{"a", 1, 1},
		{"b", 2, 2},
		{"c", 3, 3},
		{"d", 3, 4}, // 相同值，按键排序 Same value, sorted by key
		{"e", 4, 5},
	}

	for _, data := range testData {
		sl.Set(data.key, data.value)
	}

	// 测试每个元素的排名
	// Test rank of each element
	for _, data := range testData {
		rank, exists := sl.Rank(data.key)
		r.True(exists, "Key %s should exist", data.key)
		r.Equal(data.rank, rank, "Key %s: unexpected rank", data.key)
	}

	// 测试不存在的键的排名
	// Test rank of non-existent key
	rank, exists := sl.Rank("x")
	r.False(exists, "Rank should return false for non-existent key")
	r.Equal(0, rank, "Rank should be 0 for non-existent key")
}

// 测试边界情况
// Test edge cases
func TestEdgeCases(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()

	// 测试空跳表的操作
	// Test operations on empty skip list
	_, exists := sl.Get(1)
	r.False(exists, "Get should return false for empty skip list")
	r.False(sl.Del(1), "Delete should return false for empty skip list")

	// 测试大量数据
	// Test large amount of data
	for i := 0; i < 1000; i++ {
		sl.Set(i, i)
	}

	r.Equal(1000, sl.length, "Length should be 1000 after insertions")

	// 测试全部删除
	// Test deleting all elements
	for i := 0; i < 1000; i++ {
		r.True(sl.Del(i), "Failed to delete key %d", i)
	}

	r.Equal(0, sl.length, "Length should be 0 after all deletions")
}

// 测试并发操作的正确性
// Test correctness of concurrent operations
func TestConcurrent(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()
	done := make(chan bool)

	// 并发插入
	// Concurrent insertion
	for i := 0; i < 10; i++ {
		go func(val int) {
			sl.Set(val, val*10)
			done <- true
		}(i)
	}

	// 等待所有协程完成
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证结果
	// Verify results
	r.Equal(10, sl.length, "Length should be 10 after concurrent insertions")

	for i := 0; i < 10; i++ {
		val, exists := sl.Get(i)
		r.True(exists, "Key %d should exist", i)
		r.Equal(i*10, val, "Unexpected value for key %d", i)
	}
}

// 测试特殊情况下的排名
// Test ranking in special cases
func TestSpecialRanking(t *testing.T) {
	r := require.New(t)
	sl := New[int, int]()

	// 测试相同值的排名
	// Test ranking with same values
	sl.Set(1, 100)
	sl.Set(2, 100)
	sl.Set(3, 100)

	rank1, _ := sl.Rank(1)
	rank2, _ := sl.Rank(2)
	rank3, _ := sl.Rank(3)

	r.Equal(1, rank1, "First item with same value should have rank 1")
	r.Equal(2, rank2, "Second item with same value should have rank 2")
	r.Equal(3, rank3, "Third item with same value should have rank 3")
}
