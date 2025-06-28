package utils

import "math/rand"

func RandomInt() int {
	minValue := 1
	maxValue := 7
	return minValue + rand.Intn(maxValue-minValue)
}
