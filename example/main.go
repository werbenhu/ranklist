package main

import (
	"log/slog"

	"github.com/werbenhu/ranklist"
)

func main() {
	sl := ranklist.New[string, int]()

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
		// {"e", 4, 5},
	}

	for _, data := range testData {
		sl.Set(data.key, data.value)
	}

	// for _, data := range testData {
	// 	rank, exists := sl.Rank(data.key)
	// 	slog.Info("", "key", data.key, "rank", rank, "exists", exists)
	// }

	rank, exists := sl.Rank("d")
	slog.Info("", "key", "d", "rank", rank, "exists", exists)
}
