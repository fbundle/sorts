package sorts5

import (
	"github.com/fbundle/sorts/form"
)

type Form = form.Form
type Name = form.Name
type List = form.List

type Sort interface {
	Form() Form
	Level() int
	Parent() Sort
	LessEqual(dst Sort) bool
}

type Mode string

const (
	ModeComp  Mode = "COMP"  // type checking
	ModeEval  Mode = "EVAL"  // type checking and evaluation
	ModeDebug Mode = "DEBUG" // type checking and print everything
)

// Frame -
type Frame interface {
	Get(name string) Sort
	Set(name string, sort Sort) Context
	Del(name string) Context
}

type Context interface {
	Frame
}
