package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
)

func (r Context) Compile(node form.Form) (almost_sort_extra.Context, almost_sort.AlmostSort) {
	switch node := node.(type) {
	case form.Name:
		return r, r.Get(node)
	case form.List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(form.Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := r.listParsers.Get(head); ok {
			return listParser(r, node)
		} else {
			// by default, compile as beta reduction (function call)
			return almost_sort_extra.ListCompileBeta(r, node)
		}
	default:
		panic("parse_error")
	}
}
