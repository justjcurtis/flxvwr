package utils

import "math/rand/v2"

func PlaylistToString(playlist []string) string {
	str := ""
	for _, p := range playlist {
		str += p + "\n"
	}
	return str
}

func SortStrings(arr []string) {
	for i := 0; i < len(arr); i++ {
		for j := i; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

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
