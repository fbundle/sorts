package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func toListParser(listParse sorts.ListParseFunc) el_almost_sort.ListParseFunc {
	return func(parse el_almost_sort.ParseFunc, list form.List) el_almost_sort.AlmostSort {
		sort := listParse(func(form sorts.Form) sorts.Sort {
			s := el_almost_sort.MustSort(parse(form)) // inside a sort, must be sort
			if s == nil {
				panic(TypeErr)
			}
			return s
		}, list)
		return el_almost_sort.ActualSort{Sort: sort}
	}
}
