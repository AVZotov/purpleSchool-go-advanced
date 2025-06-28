package main

import (
	"fmt"
	"http_server/utils"
)

func main() {
	fmt.Println("Starting server")
	for i := 0; i < 30; i++ {
		fmt.Println(utils.RandomInt())
	}
}
