package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func (ctx Context) Form(s any) sorts.Form {
	return ctx.universe.Form(s)
}

func (ctx Context) Level(s sorts.Sort) int {
	return ctx.universe.Level(s)
}

func (ctx Context) Parent(s sorts.Sort) sorts.Sort {
	return ctx.universe.Parent(s)
}

func (ctx Context) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return ctx.universe.LessEqual(x, y)
}

func (ctx Context) LessEqualAtom(src sorts.Name, dst sorts.Name) bool {
	return ctx.universe.LessEqualAtom(src, dst)
}

func (ctx Context) NewTerm(name form.Form, parent almost_sort_extra.ActualSort) almost_sort_extra.ActualSort {
	sort := ctx.universe.NewTerm(name, parent.Repr())
	return almost_sort_extra.NewActualSort(sort)
}
