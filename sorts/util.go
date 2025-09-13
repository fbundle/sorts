package sorts

func must(a SortAttr) mustSortAttr {
	return mustSortAttr{a}
}

type mustSortAttr struct {
	a SortAttr
}

func (m mustSortAttr) lessEqual(x Sort, y Sort) {
	if !m.a.LessEqual(x, y) {
		panic(TypeErr)
	}
}

func (m mustSortAttr) termOf(x Sort, X Sort) {
	if !m.a.LessEqual(m.a.Parent(x), X) {
		panic(TypeErr)
	}
}
