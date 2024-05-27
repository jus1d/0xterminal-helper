package slice

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
