package almost_sort_extra

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Sort = sorts.Sort
type Atom = sorts.Atom

// Compiler - recursive compilation
type Compiler interface {
	Compile(form form.Form) Sort
}

// Frame - name binding
type Frame interface {
	Get(name form.Name) Sort
	Set(name form.Name, sort Sort) Context
	Del(name form.Name) Context
}

// Universe - type/el_sorts universe
type Universe interface {
	sorts.SortAttr
	Initial(level int) Sort
	Terminal(level int) Sort
	NewTerm(form form.Form, parent Sort) Atom
}

type Context interface {
	Universe
	Compiler
	Frame
}
