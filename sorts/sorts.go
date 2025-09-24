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

	// Compile - Level, Parent, LessEqual only available after compilation is done
	Compile(ctx Context) Sort

	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool
}

type Context interface {
	LessEqual(s Form, d Form) bool
}

var _ = []Sort{
	Atom{}, Type{}, Pi{},
}
