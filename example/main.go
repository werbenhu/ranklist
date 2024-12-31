package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/werbenhu/ranklist"
)

func main() {
	sl := ranklist.New[string, int]()

	for i := 0; i < 20; i++ {
		k := rand.IntN(20)
		sl.Set(strconv.Itoa(k), k)
		// sl.Set(strconv.Itoa(i), i)
		sl.Print()
	}

	for i := 0; i < 20; i++ {

		r, ok := sl.Rank(strconv.Itoa(i))

		if ok {
			fmt.Printf("rank of %d: %d\n", i, r)
		}

	}

	// sl.Del(strconv.Itoa(1))
	// sl.Print()
}
