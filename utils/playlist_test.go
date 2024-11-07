package utils

import (
	"reflect"
	"testing"
)

func slicesHaveSameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	elementCount := make(map[T]int)
	for _, elem := range a {
		elementCount[elem]++
	}
	for _, elem := range b {
		if elementCount[elem] == 0 {
			return false
		}
		elementCount[elem]--
	}
	for _, count := range elementCount {
		if count != 0 {
			return false
		}
	}
	return true
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
	}{
		{"Shuffle Empty Array", []int{}, 1},
		{"Shuffle Single Element", []int{1}, 1},
		{"Shuffle Multiple Elements", []int{1, 2, 3, 4, 5}, 3},
		{"Shuffle With Zero Rounds", []int{1, 2, 3, 4, 5}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]int, len(tt.input))
			copy(original, tt.input)
			result := Shuffle(tt.input, tt.n)

			if len(result) != len(original) {
				t.Errorf("expected length %d, got %d", len(original), len(result))
			}

			if !slicesHaveSameElements(original, result) {
				t.Errorf("shuffle altered elements; expected %v, got %v", original, result)
			}

			if len(original) > 1 && tt.n > 0 && reflect.DeepEqual(result, original) {
				t.Error("shuffle did not change the order of elements")
			}
		})
	}
}
