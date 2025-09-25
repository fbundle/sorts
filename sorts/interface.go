package sorts

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type Name = form.Name
type List = form.List
type Form = form.Form

var TypeError = errors.New("type_error")

// Sort - (or value)
type Sort interface {
	Form() Form
	Parent(ctx Context) Sort
	Level(ctx Context) int
	LessEqual(ctx Context, d Sort) bool
}

// Code - can be evaluated into Sort
type Code interface {
	Form() Form
	Eval(ctx Context) Sort
}

var _ = []Sort{
	Atom{}, Pi{}, // Inductive
}

var _ = []Code{
	Var{}, Inhabited{}, Type{}, Beta{}, Let{}, // Let, Match, etc
}

type Frame interface {
	Get(name Name) Sort
	Set(name Name, sort Sort) Context
}

type Parser interface {
	Parse(form Form) Code
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Parser
	Universe
}

type ListParseFunc = func(ctx Context, list List) Code
type NameParseFunc = func(ctx Context, name Name) Code

var FinalListParseFunc ListParseFunc
var FinalNameParseFunc NameParseFunc

var ListParseFuncMap = map[Name]ListParseFunc{}
