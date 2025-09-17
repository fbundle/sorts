package sorts5

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

type Frame interface {
	Get(name string) Sort
	Set(name string, sort Sort) Context
	Del(name string) Context
}

type ListParser struct {
	Command   Name
	ListParse func(ctx Context, list List) (Context, Sort)
}

type Parser interface {
	Parse(form Form) Sort
	AddListParser(listParser ListParser) Context
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Parser
	Universe
}
