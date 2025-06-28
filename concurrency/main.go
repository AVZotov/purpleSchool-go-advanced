package main

import (
	"fmt"
	"math/rand"
)

func main() {
	in := make(chan int)
	out := make(chan int)
	sl := getSlice()
	for _, v := range sl {
		go func() {
			out <- v
		}()
	}

	for range out {
		go func() {
			in <- power(<-out)
		}()
	}

	for range in {
		fmt.Printf("value is: %d\n", <-in)
	}
}

func getSlice() []int {
	sl := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		val := rand.Intn(101)
		sl = append(sl, val)
	}
	return sl
}

func power(value int) int {
	return value * value
}
