package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
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

func (ctx Context) GetRule(src sorts.Name, dst sorts.Name) bool {
	return ctx.universe.GetRule(src, dst)
}

func (ctx Context) NewTerm(name form.Name, parent almost_sort.ActualSort) almost_sort.ActualSort {
	sort := ctx.universe.NewTerm(name, parent.Repr())
	return almost_sort.NewActualSort(sort)
}
