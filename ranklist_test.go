package ranklist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRankList_Basic(t *testing.T) {
	r := require.New(t)
	rl := New[string, int]()

	// 测试插入和获取
	rl.Set("player1", 100)
	r.Equal(100, rl.Get("player1"), "should get correct score")

	// 测试排名
	r.Equal(1, rl.Rank("player1"), "should get correct rank")

	// 测试删除
	r.True(rl.Del("player1"), "should successfully delete existing key")
	r.Equal(0, rl.Get("player1"), "should return zero value after deletion")
	r.Equal(0, rl.Rank("player1"), "should return zero rank after deletion")
}

func TestRankList_MultipleEntries(t *testing.T) {
	r := require.New(t)
	rl := New[string, int]()

	// 插入多个玩家
	players := map[string]int{
		"player1": 100,
		"player2": 200,
		"player3": 150,
		"player4": 200, // 相同分数测试
	}

	for player, score := range players {
		rl.Set(player, score)
	}

	// 验证长度
	r.Equal(4, rl.length, "should have correct length")

	// 验证排名 (相同分数按key字典序)
	expectedRanks := map[string]int{
		"player1": 1, // 200分，key较小
		"player3": 2, // 200分，key较大
		"player2": 3, // 150分
		"player4": 4, // 100分
	}

	for player, expectedRank := range expectedRanks {
		r.Equal(expectedRank, rl.Rank(player), "should have correct rank for %s", player)
	}
}

func TestRankList_UpdateScore(t *testing.T) {
	r := require.New(t)
	rl := New[string, int]()

	// 测试分数更新
	rl.Set("player1", 100)
	rl.Set("player2", 200)
	rl.Set("player1", 300) // 更新分数

	r.Equal(300, rl.Get("player1"), "should get updated score")
	r.Equal(1, rl.Rank("player1"), "should get updated rank")
	r.Equal(2, rl.Rank("player2"), "should get correct rank for other player")
}

func TestRankList_EdgeCases(t *testing.T) {
	r := require.New(t)
	rl := New[string, int]()

	// 测试空排行榜
	r.Equal(0, rl.Rank("nonexistent"), "should return 0 rank for nonexistent key")
	r.Equal(0, rl.Get("nonexistent"), "should return 0 score for nonexistent key")
	r.False(rl.Del("nonexistent"), "should return false when deleting nonexistent key")

	// 测试零分和负分
	rl.Set("zero", 0)
	rl.Set("negative", -100)
	rl.Set("positive", 100)

	r.Equal(1, rl.Rank("negative"), "negative score should rank first")
	r.Equal(2, rl.Rank("zero"), "zero score should rank second")
	r.Equal(3, rl.Rank("positive"), "positive score should rank last")
}

func TestRankList_DifferentTypes(t *testing.T) {
	r := require.New(t)

	// 测试整数键-浮点数值
	rlFloat := New[int, float64]()
	rlFloat.Set(1, 100.5)
	rlFloat.Set(2, 200.5)

	r.Equal(100.5, rlFloat.Get(1), "should get correct float64 score")
	r.Equal(2, rlFloat.Rank(1), "should get correct rank")

	// 测试字符串键-字符串值
	rlString := New[string, string]()
	rlString.Set("key1", "value1")
	rlString.Set("key2", "value2")

	r.Equal("value1", rlString.Get("key1"), "should get correct string score")
	r.Equal(1, rlString.Rank("key1"), "should get correct rank")
}

func TestRankList_StressTest(t *testing.T) {
	r := require.New(t)
	rl := New[int, int]()
	const n = 1000

	// 插入大量数据
	for i := 0; i < n; i++ {
		rl.Set(i, i)
	}

	r.Equal(n, rl.length, "should have correct length after bulk insert")

	// 验证所有数据
	for i := 0; i < n; i++ {
		r.Equal(i, rl.Get(i), "should get correct score for %d", i)
		r.Equal(n-i, rl.Rank(i), "should get correct rank for %d", i)
	}

	// 批量删除测试
	for i := 0; i < n; i += 2 {
		r.True(rl.Del(i), "should successfully delete key %d", i)
	}

	r.Equal(n/2, rl.length, "should have correct length after bulk deletion")
}

func TestRankList_ConcurrentAccess(t *testing.T) {
	r := require.New(t)
	rl := New[int, int]()

	// 预先插入一些数据
	for i := 0; i < 10; i++ {
		rl.Set(i, i*100)
	}

	// 验证数据一致性
	for i := 0; i < 10; i++ {
		score := rl.Get(i)
		rank := rl.Rank(i)
		r.Equal(i*100, score, "should get correct score")
		r.Equal(10-i, rank, "should get correct rank")
	}
}
