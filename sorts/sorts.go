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

	// Compile - return an atom of the correct type
	Compile(ctx Context) Sort // Level, Parent, LessEqual, Reduce only available after compilation is done
	Level(ctx Context) int
	Parent(ctx Context) Sort
	LessEqual(ctx Context, d Sort) bool

	// Reduce - return the actual output
	Reduce(ctx Context) Sort
}

type Context interface {
	Set(name Name, sort Sort) Context
	LessEqual(s Form, d Form) bool
}

var _ = []Sort{
	Atom{}, Type{}, Pi{}, Beta{},
}
