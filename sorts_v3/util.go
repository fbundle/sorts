package sorts

func mustTermOf(x Sort, X Sort) {
	mustSubType(Parent(x), X)
}
func mustSubType(x Sort, y Sort) {
	if ok := LessEqual(x, y); !ok {
		panic("type_error")
	}
}
