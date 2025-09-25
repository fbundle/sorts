package sorts

import (
	"github.com/fbundle/sorts/form"
)

type Name = form.Name
type List = form.List
type Form = form.Form

type Sort interface {
	Form() Form

	Parent(ctx Context) Sort // type checking
	Level(ctx Context) int
	LessEqual(ctx Context, d Sort) bool

	Eval(ctx Context) Sort // evaluation
}

type Code interface {
	Form() Form
}

var _ = []Sort{
	Atom{}, Lambda{}, Beta{}, Type{}, Inhabited{}, Let{},
}

var _ = []Code{
	Annot{}, Binding{},
}

type Frame interface {
	Set(name Name, sort Sort) Context
}
type ListParseFunc = func(ctx Context, list List) Sort

type Parser interface {
	Parse(form Form) Sort
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Mode string

const (
	ModeComp  Mode = "COMP"  // type checking
	ModeEval  Mode = "EVAL"  // type checking and evaluation
	ModeDebug Mode = "DEBUG" // type checking and print everything
)

type Context interface {
	Frame
	Parser
	Universe
	Mode() Mode
}

var DefaultParseFunc ListParseFunc

var ListParseFuncMap = map[Name]ListParseFunc{}
