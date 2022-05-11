package iu

func SliceContains[T comparable](vs []T, t T) bool {
	for _, v := range vs {
		if v == t {
			return true
		}
	}
	return false
}
