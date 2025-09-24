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
	Level(ctx Context) int
	LessEqual(ctx Context, d Sort) bool
}

type Context interface {
	LessEqual(s Form, d Form) bool
}
