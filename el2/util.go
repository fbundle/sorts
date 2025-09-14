package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func toListParser(listParse sorts.ListParseFunc) ListParseFunc {
	return func(parse ParseFunc, list List) el_almost_sort.AlmostSort {
		sort := listParse(func(form sorts.Form) sorts.Sort {
			s := el_almost_sort.mustSort(parse(form)) // inside a sort, must be sort
			if s == nil {
				panic(TypeErr)
			}
			return s
		}, list)
		return el_almost_sort.ActualSort{sort}
	}
}
