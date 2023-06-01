package utils

func Includes[T comparable](slice []T, findItem T) int {
	for index, item := range slice {
		if item == findItem {
			return index
		}
	}
	return -1
}
