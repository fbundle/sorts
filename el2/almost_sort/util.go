package el_almost_sort

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func must(a sorts.SortAttr) mustSortAttr {
	return mustSortAttr{a}
}

type mustSortAttr struct {
	a sorts.SortAttr
}

func (m mustSortAttr) lessEqual(x sorts.Sort, y sorts.Sort) {
	if !m.a.LessEqual(x, y) {
		panic(TypeErr)
	}
}

func (m mustSortAttr) termOf(x sorts.Sort, X sorts.Sort) {
	if !m.a.LessEqual(m.a.Parent(x), X) {
		panic(TypeErr)
	}
}

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
		panic(TypeErr)
	}
}
