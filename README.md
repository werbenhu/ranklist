# ranklist

一个基于跳表(Skip List)实现的高性能排行榜数据结构。
A high-performance ranking system implementation based on Skip List.

## 特性 | Features

- 支持泛型，可用于各种可比较的数据类型
- O(log n) 的时间复杂度用于插入、删除和查询操作
- 维护实时排名
- 支持相同分数下的二级排序
- 内置哈希表实现快速查找

---

- Generic support for all comparable types
- O(log n) time complexity for insertion, deletion and query operations
- Maintains real-time rankings
- Supports secondary sorting for equal scores
- Built-in hash table for quick lookups

## 使用示例 | Usage

```go
// 创建新的排行榜
rankList := ranklist.New[string, int]()

// 插入数据
rankList.Set("player1", 100)
rankList.Set("player2", 200)
rankList.Set("player3", 150)

// 获取分数
score := rankList.Get("player1")  // 返回 100

// 获取排名
rank := rankList.Rank("player2")  // 返回 1 (最高分)

// 删除数据
rankList.Del("player3")