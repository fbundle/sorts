package sort
/*
the code below was auto generated for inductive type
inductive Sort
	| Atom	: string -> int -> Sort -> Sort
	| Arrow	: Sort -> Sort -> Sort
	| Prod	: Sort -> Sort -> Sort
	| Sum	: Sort -> Sort -> Sort
*/
type Sort[T any] interface {
	attrSort()
}
type Atom[T any] struct {
	Field_0 string
	Field_1 int
	Field_2 Sort[T]
}
func (o Atom[T]) attrSort() {}
func (o Atom[T]) Unwrap() (string,int,Sort[T]) {
	return o.Field_0 , o.Field_1 , o.Field_2
}
type Arrow[T any] struct {
	Field_0 Sort[T]
	Field_1 Sort[T]
}
func (o Arrow[T]) attrSort() {}
func (o Arrow[T]) Unwrap() (Sort[T],Sort[T]) {
	return o.Field_0 , o.Field_1
}
type Prod[T any] struct {
	Field_0 Sort[T]
	Field_1 Sort[T]
}
func (o Prod[T]) attrSort() {}
func (o Prod[T]) Unwrap() (Sort[T],Sort[T]) {
	return o.Field_0 , o.Field_1
}
type Sum[T any] struct {
	Field_0 Sort[T]
	Field_1 Sort[T]
}
func (o Sum[T]) attrSort() {}
func (o Sum[T]) Unwrap() (Sort[T],Sort[T]) {
	return o.Field_0 , o.Field_1
}
type Match[T any, V any] struct {
	MapAtom func(Atom[T]) V
	MapArrow func(Arrow[T]) V
	MapProd func(Prod[T]) V
	MapSum func(Sum[T]) V
}
func (m Match[T, V]) Apply(o Sort[T]) V {
	switch o := o.(type) {
		case Atom[T]:
			return m.MapAtom(o)
		case Arrow[T]:
			return m.MapArrow(o)
		case Prod[T]:
			return m.MapProd(o)
		case Sum[T]:
			return m.MapSum(o)
		default:
			panic("unreachable")
	}
}