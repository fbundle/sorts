package sorts

type PartialOrder interface {
	LessEqual(dst PartialOrder) bool
}

type Sort[T PartialOrder] interface {
	sortAttr() sortAttr[T]
}

type sortAttr[T PartialOrder] struct {
	data      T
	level     int
	parent    Sort[T]
	lessEqual func(dst Sort[T]) bool
}

func Data[T PartialOrder](s Sort[T]) T {
	return s.sortAttr().data
}
func Level[T PartialOrder](s Sort[T]) int {
	return s.sortAttr().level
}
func Parent[T PartialOrder](s Sort[T]) Sort[T] {
	return s.sortAttr().parent
}
func LessEqual[T PartialOrder](x Sort[T], y Sort[T]) bool {
	return x.sortAttr().lessEqual(y)
}
