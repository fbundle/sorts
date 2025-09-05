package sorts

func mustType(x Sort, X Sort) {
	mustEqual(Parent(x), X)
}

func mustEqual(x Sort, y Sort) {
	if x != y {
		panic("type_error")
	}
}
