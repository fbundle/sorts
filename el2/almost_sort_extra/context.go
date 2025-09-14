package almost_sort_extra

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

// Compiler - recursive compilation
type Compiler interface {
	Compile(form form.Form) (Context, AlmostSort)
}

// Frame - name binding
type Frame interface {
	Get(name form.Name) ActualSort
	Set(name form.Name, sort ActualSort) Context
	Del(name form.Name) Context
}

// Universe - type/sort universe
type Universe interface {
	sorts.SortAttr
	NewTerm(name form.Form, parent ActualSort) ActualSort
}

type Context interface {
	Universe
	Compiler
	Frame
}
