package el2_almost_sort

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Compiler interface {
	Compile(form form.Form) (Context, almost_sort.AlmostSort)
}
type Frame interface {
	Get(name form.Name) almost_sort.ActualSort
	Set(name form.Name, sort almost_sort.ActualSort) Context
}

type Context interface {
	sorts.SortAttr
	Compiler
	Frame
}

type ListParseFunc = func(r Context, list form.List) (Context, almost_sort.AlmostSort)
