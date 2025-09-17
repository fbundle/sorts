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

type Binding struct {
	Name Name
	Sort Sort
}

type Frame interface {
	Get(name Name) Sort
	Set(name Name, sort Sort) Context
}
type ListParseFunc = func(ctx Context, list List) (Sort, []Binding)

type Parser interface {
	Parse(form Form) (Sort, []Binding)
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Parser
	Universe
}

var ListParseFuncMap = map[Name]ListParseFunc{}
