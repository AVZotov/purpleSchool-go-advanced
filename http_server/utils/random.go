package utils

import (
	"math/rand"
	"time"
)

func RandomInt() int {
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	minValue := 1
	maxValue := 7
	return minValue + rnd.Intn(maxValue-minValue)
}
