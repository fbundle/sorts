package sorts

func mustTermOf(x WithSort, X WithSort) {
	mustSubType(Parent(x), X)
}
func mustSubType(x WithSort, y WithSort) {
	if ok := LessEqual(x, y); !ok {
		panic("type_error")
	}
}
func dummyTerm(parent WithSort, name string) WithSort {
	return NewAtom(Level(parent)-1, name, parent)
}
