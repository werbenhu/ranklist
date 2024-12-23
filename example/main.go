package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/werbenhu/skiplist"
)

func main() {
	sl := skiplist.NewSkiplist()

	for i := 0; i < 1000; i++ {
		sl.Insert(i, float64(i))

	}

	for i := 0; i < 10; i++ {
		k := rand.IntN(1000)
		score := sl.GetScore(k)
		rank := sl.GetRank(k)
		fmt.Println(k, score, rank)
	}

}
