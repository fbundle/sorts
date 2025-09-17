package sorts

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

/*
stage 1 (parsing):		Form -> Sort
stage 2 (compiling):	Sort -> Sort
stage 3 (reducing):		Sort -> Sort
*/

type Sort interface {
	Form() Form

	Compile(ctx Context) Sort
	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool

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
