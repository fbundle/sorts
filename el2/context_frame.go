package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
)

func (ctx Context) Get(name form.Name) almost_sort.ActualSort {
	if s, ok := ctx.frame.Get(name); ok {
		return s
	}
	if s, ok := ctx.universe.GetBuiltin(name); ok {
		return almost_sort.NewActualSort(s)
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (ctx Context) Set(name form.Name, sort almost_sort.ActualSort) almost_sort_extra.Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}
