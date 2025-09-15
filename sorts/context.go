package sorts

import (
	"github.com/fbundle/sorts/form"
)

type Mode string

const (
	ModeComp  Mode = "COMP"  // type checking
	ModeEval  Mode = "EVAL"  // type checking and evaluation
	ModeDebug Mode = "DEBUG" // type checking and print everything
)

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
	SortAttr
	Initial(level int) Sort
	Terminal(level int) Sort
	NewTerm(form form.Form, parent Sort) Atom
}

type Context interface {
	Universe
	Compiler
	Frame
	Mode() Mode
	ToString(o any) string
}
type ListCompileFunc = func(r Context, list form.List) Sort
