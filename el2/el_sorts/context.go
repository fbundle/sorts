package almost_sort_extra

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

// Compiler - recursive compilation
type Compiler interface {
	Compile(form form.Form) (Context, Sort)
}

// Frame - name binding
type Frame interface {
	Get(name form.Name) typeSort
	Set(name form.Name, sort typeSort) Context
	Del(name form.Name) Context
}

// Universe - type/el_sorts universe
type Universe interface {
	sorts.SortAttr
	NewTerm(name form.Form, parent typeSort) typeSort
}

type Context interface {
	Universe
	Compiler
	Frame
}
