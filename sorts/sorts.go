package sorts

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type Name = form.Name
type List = form.List
type Form = form.Form

var TypeErr = errors.New("type_error")

/*
stage 1 (parsing):		Form -> Sort
stage 2 (compiling):	Sort -> Sort
stage 3 (reducing):		Sort -> Sort
*/

type Sort interface {
	Form() Form

	TypeCheck(ctx Context) Sort

	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool

	Reduce(ctx Context) Sort
}

type Frame interface {
	Get(name Name) Sort
	Set(name Name, sort Sort) Context
}
type ListCompileFunc = func(ctx Context, list List) Sort

type Compiler interface {
	Compile(form Form) Sort
}

type Universe interface {
	LessEqual(src Form, dst Form) bool
}

type Context interface {
	Frame
	Compiler
	Universe
}

var DefaultCompileFunc ListCompileFunc

var ListCompileFuncMap = map[Name]ListCompileFunc{}
