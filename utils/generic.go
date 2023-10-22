package utils

func CopySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func CopyStruct[T any](src *T) *T {
	dst := new(T)
	*dst = *src
	return dst
}
