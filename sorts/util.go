package sorts

func must(a SortAttr) mustSortAttr {
	return mustSortAttr{a}
}

type mustSortAttr struct {
	SortAttr
}

func (m mustSortAttr) mustLessEqual(x Sort, y Sort) {
	if !m.LessEqual(x, y) {
		panic(TypeErr)
	}
}

func (m mustSortAttr) mustTermOf(x Sort, X Sort) {
	if !m.LessEqual(m.Parent(x), X) {
		panic(TypeErr)
	}
}
