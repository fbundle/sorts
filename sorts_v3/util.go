package sorts

func MustTermOf(x Sort, X Sort) {
	X1 := Parent(x)
	if ok := LessEqual(X1, X); !ok {
		panic("type_error")
	}
}
