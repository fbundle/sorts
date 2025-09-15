package list

/*
the code below was auto generated for inductive type
inductive List

	| Nil	: List
	| Cons	: T -> List -> List
*/
type List[T any] interface {
	attrList()
}
type Nil[T any] struct {
}

func (o Nil[T]) attrList() {}
func (o Nil[T]) Unwrap() {
	return
}

type Cons[T any] struct {
	Field_0 T
	Field_1 List[T]
}

func (o Cons[T]) attrList() {}
func (o Cons[T]) Unwrap() (T, List[T]) {
	return o.Field_0, o.Field_1
}

type Match[T any, V any] struct {
	MapNil  func(Nil[T]) V
	MapCons func(Cons[T]) V
}

func (m Match[T, V]) Apply(o List[T]) V {
	switch o := o.(type) {
	case Nil[T]:
		return m.MapNil(o)
	case Cons[T]:
		return m.MapCons(o)
	default:
		panic("unreachable")
	}
}
