package sorts

func MustTermOf(x Sort, X Sort) {
	X1 := Parent(x)
	if ok := SubTypeOf(X1, X); !ok {
		panic("type_error")
	}
}

func TermOf(x Sort, X Sort) bool {
	X1 := Parent(x)
	return SubTypeOf(X1, X)
}
