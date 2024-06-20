package pointer

func Ptr[T any](v T) *T {
	return &v
}

func GetValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
