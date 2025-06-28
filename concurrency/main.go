package main

import (
	"fmt"
	"math/rand"
)

func main() {
	in := make(chan int)
	out := make(chan int)

	sendValues(in)
	processValues(in, out)

	for v := range out {
		fmt.Printf("Result: %d\n", v)
	}
}

func sendValues(ch chan<- int) {
	go func() {
		defer close(ch)
		sl := getSlice()
		for _, v := range sl {
			ch <- v
		}
	}()
}

func processValues(chIn <-chan int, chOut chan<- int) {
	go func() {
		defer close(chOut)
		for value := range chIn {
			square := power(value)
			chOut <- square
		}
	}()
}

func power(value int) int {
	return value * value
}

func getSlice() []int {
	sl := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		val := rand.Intn(101)
		sl = append(sl, val)
	}
	return sl
}
