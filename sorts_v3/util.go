package sorts

func mustType(x Sort, X Sort) {
	mustLessEqual(Parent(x), X)
}
func mustLessEqual(x Sort, y Sort) {
	if ok := LessEqual(x, y); !ok {
		panic("type_error")
	}
}
