package sorts

func mustTermOf(x Sort, X Sort) {
	X1 := Parent(x)
	if ok := LessEqual(X1, X); !ok {
		panic("type_error")
	}
}

func dummyTerm(parent Sort, name string) Sort {
	return newAtom(Level(parent)-1, name, parent)
}
