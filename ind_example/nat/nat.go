package nat
/*
the code below was auto generated for inductive type
inductive Nat
	| Zero	: Nat
	| Succ	: Nat -> Nat
*/
type Nat[T any] interface {
	attrNat()
}
type Zero[T any] struct {
}
func (o Zero[T]) attrNat() {}
func (o Zero[T]) Unwrap() () {
	return 
}
type Succ[T any] struct {
	Field_0 Nat[T]
}
func (o Succ[T]) attrNat() {}
func (o Succ[T]) Unwrap() (Nat[T]) {
	return o.Field_0
}
type Match[T any, V any] struct {
	MapZero func(Zero[T]) V
	MapSucc func(Succ[T]) V
}
func (m Match[T, V]) Apply(o Nat[T]) V {
	switch o := o.(type) {
		case Zero[T]:
			return m.MapZero(o)
		case Succ[T]:
			return m.MapSucc(o)
		default:
			panic("unreachable")
	}
}