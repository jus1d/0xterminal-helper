package slice

import (
	"math/rand"
	"time"
)

func Unique[T comparable](slice []T) []T {
	unique := make([]T, 0, len(slice))
	seen := make(map[T]struct{})

	for _, v := range slice {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		unique = append(unique, v)
	}

	return unique
}

func Choose[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(slice))
	return slice[index]
}
