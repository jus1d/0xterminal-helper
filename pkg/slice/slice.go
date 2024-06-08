package slice

import (
	"math/rand"
	"time"
)

// Unique returns a new slice with unique elements.
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

// Choose returns a random element from the provided slice.
func Choose[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	index := r.Intn(len(slice))
	return slice[index]
}
