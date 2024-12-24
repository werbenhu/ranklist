package main

import (
	"fmt"

	"github.com/werbenhu/ranklist"
)

func main() {
	sl := ranklist.New[int, int]()

	for i := 0; i < 10; i++ {
		sl.Set(i, i)
	}

	for i := 0; i < 10; i++ {
		sl.Set(i+10, i)
	}
	// sl.Print()

	for i := 0; i < 10; i++ {
		// k := rand.IntN(1000000)
		k := i + 10
		score := sl.Get(k)
		rank := sl.Get(k)
		fmt.Println(k, score, rank)
	}
}
