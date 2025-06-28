package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	in := make(chan int)
	out := make(chan int)
	sl := getSlice()
	var wg1, wg2 sync.WaitGroup

	wg1.Add(len(sl))
	for _, v := range sl {
		go func(value int) {
			defer wg1.Done()
			out <- value
		}(v)
	}

	go func() {
		wg1.Wait()
		close(out)
	}()

	wg2.Add(len(sl))
	for range out {
		go func() {
			defer wg2.Done()
			value := <-in
			squared := power(value)
			fmt.Printf("Обработка: %d^2 = %d\n", value, squared)
		}()
	}

	go func() {
		wg2.Wait()
		close(in)
	}()

	fmt.Println("\nFinal Results:")
	for value := range in {
		fmt.Printf("value is: %d\n", value)
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
