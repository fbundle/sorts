package el2

import (
	"cmp"
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Name = form.Name
type Form = form.Form
type List = form.List
type Sort = sorts.Sort
type SortAttr = sorts.SortAttr

var NewAtomChain = sorts.NewAtomChain
var NewAtomTerm = sorts.NewAtomTerm

type ParseFunc = func(form Form) AlmostSort
type ListParseFunc = func(parse ParseFunc, list List) AlmostSort
type ParseSortFunc = sorts.ParseFunc
type ListParseSortFunc = sorts.ListParseFunc

var ListParseArrow = sorts.ListParseArrow
var ListParseSum = sorts.ListParseSum
var ListParseProd = sorts.ListParseProd

type rule struct {
	src sorts.Name
	dst sorts.Name
}

func (r rule) Cmp(s rule) int {
	if c := cmp.Compare(r.src, s.src); c != 0 {
		return c
	}
	return cmp.Compare(r.dst, s.dst)
}

func must(a SortAttr) mustSortAttr {
	return mustSortAttr{a}
}

var TypeErr = fmt.Errorf("type_error")

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
