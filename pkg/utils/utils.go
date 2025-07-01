package utils

// Contains returns true if the slice contains the element; false otherwise.
func Contains[T comparable](slice []T, elem T) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// GetValueOrDefault returns the value pointed to by p if it is not nil.
// If p is nil and a default value is provided, it returns that default.
// Otherwise, it returns the zero value for type T.
func GetValueOrDefault[T any](p *T, defaultValue ...T) T {
	if p != nil {
		return *p
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	var zero T
	return zero
}
