package sorts

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

type Sort1 interface {
	Compile(ctx Context) Sort2
	Form() Form
}

type Sort2 interface {
	Sort1
	Reduce(ctx Context) Sort3

	Level(ctx Context) int
	Parent(ctx Context) Sort1
	LessEqual(ctx Context, d Sort1) bool
}

type Sort3 interface {
	Sort1
	Sort2
}

type Frame interface {
	Get(name Name) Sort1
	Set(name Name, sort Sort1) Context
}

type Parser interface {
	Parse(form Form) (Context, Sort1)
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Parser
	Universe
}

type ListParseFunc = func(ctx Context, list List) (Context, Sort1)

var ListParseFuncMap = map[Name]ListParseFunc{}
