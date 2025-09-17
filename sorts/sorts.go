package sorts

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

type Sort interface {
	Compile(ctx Context) Sort
	Form() Form
	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool
}

type Reducible interface {
	Sort
	Reduce(ctx Context) Sort
}

type Frame interface {
	Get(name Name) Sort
	Set(name Name, sort Sort) Context
}

type Parser interface {
	Parse(form Form) (Context, Sort)
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Parser
	Universe
}

type ListParseFunc = func(ctx Context, list List) (Context, Sort)

var ListParseFuncMap = map[Name]ListParseFunc{}
