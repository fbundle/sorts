package sorts

func mustType[T any](err error, o any) T {
	if v, ok := o.(T); ok {
		return v
	}
	panic(err)
}
