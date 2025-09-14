package el2_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Runtime interface {
	Parse(form form.Form) (Runtime, almost_sort.AlmostSort)
	SortAttr() sorts.SortAttr
	Get(name form.Name) almost_sort.ActualSort
	Set(name form.Name, sort almost_sort.ActualSort) Runtime
}

type ListParseFunc = func(r Runtime, list form.List) (Runtime, almost_sort.AlmostSort)
