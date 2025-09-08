package sorts

func mustTermOf(x Sort, X Sort) {
	mustSubType(Parent(x), X)
}
func mustSubType(x Sort, y Sort) {
	if ok := LessEqual(x, y); !ok {
		panic("type_error")
	}
}
func dummyTerm(parent Sort, name string) Sort {
	return newAtom(Level(parent)-1, name, parent)
}
