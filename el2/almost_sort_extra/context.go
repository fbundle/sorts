package almost_sort_extra

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

// Compiler - recursive compilation
type Compiler interface {
	Compile(form form.Form) (Context, almost_sort.AlmostSort)
}

// Frame - name binding
type Frame interface {
	Get(name form.Name) almost_sort.ActualSort
	Set(name form.Name, sort almost_sort.ActualSort) Context
}

// Universe - type/sort universe
type Universe interface {
	sorts.SortAttr
	NewTerm(name form.Name, parent almost_sort.ActualSort) almost_sort.ActualSort
}

type Context interface {
	Universe
	Compiler
	Frame
}
