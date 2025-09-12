package sorts

func MustTermOf(x Sort, X Sort) {
	Must(TermOf(x, X))
}

func TermOf(x Sort, X Sort) bool {
	X1 := Parent(x)
	return SubTypeOf(X1, X)
}

func Must(ok bool) {
	if !ok {
		panic("type_error")
	}
}
