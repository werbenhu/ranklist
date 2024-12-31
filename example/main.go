package main

import (
	"strconv"

	"github.com/werbenhu/ranklist"
)

func main() {
	sl := ranklist.New[string, int]()

	for i := 0; i < 20; i++ {
		// k := rand.IntN(20)
		sl.Set(strconv.Itoa(i), i)
		sl.Print()
	}

	sl.Del(strconv.Itoa(1))
	sl.Print()

	// sl.Print()
	// for i := 0; i < 3; i++ {
	// 	sl.Set(strconv.Itoa(i+10), i)
	// }
	// sl.Print()

	// for i := 0; i < 6; i++ {
	// 	// k := rand.IntN(1000000)
	// 	k := i + 10
	// 	score := sl.Get(strconv.Itoa(k))
	// 	rank := sl.Get(strconv.Itoa(k))
	// 	fmt.Println(k, score, rank)
	// }

	// rl := ranklist.New[string, int]()

	// // 测试分数更新
	// // rl.Set("player3", 150)
	// rl.Set("player1", 100)
	// rl.Set("player2", 200)
	// rl.Print()

	// rl.Set("player1", 300) // 更新分数
	// rl.Print()

	// rank := rl.Rank("player1")
	// fmt.Println(rank)
}
