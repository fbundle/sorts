package adt

type Result[T any] struct {
	Val T
	Err error
}

func (r Result[T]) Unwrap(val *T) error {
	if val != nil {
		*val = r.Val
	}
	return r.Err
}

func (r Result[T]) Iter(yield func(T)) {
	if r.Err == nil {
		yield(r.Val)
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		Err: err,
	}
}

func Ok[T any](val T) Result[T] {
	return Result[T]{
		Val: val,
		Err: nil,
	}
}
