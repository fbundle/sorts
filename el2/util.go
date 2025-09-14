package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func toAlmostSortListParser(listParse sorts.ListParseFuncWithHead) el_almost_sort.ListParseFuncWithHead {
	return func(H form.Name) el_almost_sort.ListParseFunc {
		sortListParse := listParse(H)
		return func(parse el_almost_sort.ParseFunc, list form.List) el_almost_sort.AlmostSort {
			sort := sortListParse(func(form sorts.Form) sorts.Sort {
				return el_almost_sort.MustSort(parse(form))
			}, list)
			return el_almost_sort.ActualSort{Sort: sort}
		}
	}
}
