package adt

type Iterable[T any] interface {
	Iter(yield func(T) bool)
	NonEmpty() Option[NonEmptyIterable[T]]
}

type NonEmptyIterable[T any] interface {
	Iterable[T]
	Head() T
	Tail() Iterable[T]
}
