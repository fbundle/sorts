package sorts5

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

// Frame -
type Frame interface {
	Get(name string) Sort
	Set(name string, sort Sort) Context
	Del(name string) Context
}

type Universe interface {
	Initial() Name
	Terminal() Name
	WithLessEqual(src Form, dst Form) Context
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
}

type Sort interface {
	Compile(ctx Context) Sort
	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, dst Sort) bool
}
