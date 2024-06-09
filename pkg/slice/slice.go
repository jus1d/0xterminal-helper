package slice

import (
	"crypto/rand"
	"math/big"
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
	var zero T
	if len(slice) == 0 {
		return zero
	}

	max := big.NewInt(int64(len(slice)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return zero
	}

	return slice[n.Int64()]
}
