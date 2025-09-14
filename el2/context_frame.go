package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
)

func (ctx Context) Get(name form.Name) almost_sort.ActualSort {
	s, ok := ctx.frame.Get(name)
	if !ok {
		panic(TypeErr)
	}
	return s
}

func (ctx Context) Set(name form.Name, sort almost_sort.ActualSort) almost_sort_extra.Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}
