package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/werbenhu/skiplist"
)

func main() {
	sl := skiplist.NewSkiplist[int, int]()

	for i := 0; i < 1000000; i++ {
		sl.Insert(i, i)
	}
	// sl.Print()

	for i := 0; i < 10; i++ {
		k := rand.IntN(1000000)
		score := sl.GetScore(k)
		rank := sl.GetRank(k)
		fmt.Println(k, score, rank)
	}
}
