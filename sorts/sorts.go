package sorts

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type Name = form.Name
type List = form.List
type Form = form.Form

var TypeErr = errors.New("type_error")

type Sort interface {
	Form() Form

	Parent(ctx Context) Sort // type checking
	Level(ctx Context) int
	LessEqual(ctx Context, d Sort) bool

	Reduce(ctx Context) Sort
}

type Context interface {
	Set(name Name, sort Sort) Context
	LessEqual(s Form, d Form) bool
}

var _ = []Sort{
	Atom{}, Type{}, Pi{}, Beta{},
}
