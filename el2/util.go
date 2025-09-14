package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func toAlmostSortListParser(listParse sorts.ListParseFuncWithHead) almost_sort_extra.ListParseFuncWithHead {
	return func(H form.Name) almost_sort_extra.ListCompileFunc {
		sortListParse := listParse(H)
		return func(parse almost_sort_extra.ParseFunc, list form.List) almost_sort.AlmostSort {
			sort := sortListParse(func(form sorts.Form) sorts.Sort {
				return almost_sort.MustSort(parse(form))
			}, list)
			return almost_sort.ActualSort{sort: sort}
		}
	}
}
