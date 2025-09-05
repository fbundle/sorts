package adt

type Option[T any] struct {
	Val T
	Ok  bool
}

func (o Option[T]) Unwrap(val *T) bool {
	if val != nil {
		*val = o.Val
	}
	return o.Ok
}

func (o Option[T]) MustUnwrap() T {
	if !o.Ok {
		panic(o)
	}
	return o.Val
}

func None[T any]() Option[T] {
	return Option[T]{
		Ok: false,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Val: val,
		Ok:  true,
	}
}

type Option2[T1 any, T2 any] struct {
	Val1 T1
	Val2 T2
	Ok   bool
}

func (o Option2[T1, T2]) Unwrap(val1 *T1, val2 *T2) bool {
	if val1 != nil {
		*val1 = o.Val1
	}
	if val2 != nil {
		*val2 = o.Val2
	}
	return o.Ok
}

func None2[T1 any, T2 any]() Option2[T1, T2] {
	return Option2[T1, T2]{
		Ok: false,
	}
}

func Some2[T1 any, T2 any](val1 T1, val2 T2) Option2[T1, T2] {
	return Option2[T1, T2]{
		Val1: val1,
		Val2: val2,
		Ok:   true,
	}
}
