package utils

import "math/rand/v2"

func Shuffle[T any](arr []T, n int) []T {
	for round := 0; round < n; round++ {
		for i := range arr {
			j := rand.IntN(i + 1)
			if i == j {
				continue
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	return arr

}
