package sorts

func mustType(x Sort, X Sort) {
	if ok := LessEqual(Parent(x), X); !ok {
		panic("type_error")
	}
}
