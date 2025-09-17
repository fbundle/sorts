package sorts5

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

type Frame[T any] interface {
	Get(name string) Sort
	Set(name string, sort Sort) T
	Del(name string) T
}

type ListParseFunc[T any] = func(ctx T, list List) (T, Sort)

type Parser[T any] interface {
	Parse(form Form) Sort
	AddListParseFunc(cmd Name, parseFunc ListParseFunc[T]) T
}

type Universe[T any] interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame[Context]
	Parser[Context]
	Universe[Context]
}

type Sort interface {
	Compile(ctx Context) Sort
	Form() Form
	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool
}
