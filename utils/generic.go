package utils

func ClonePointer[T any](src *T) *T {
	dst := new(T)
	*dst = *src
	return dst
}
