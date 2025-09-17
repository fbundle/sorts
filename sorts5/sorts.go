package sorts5

import (
	"github.com/fbundle/sorts/form"
)

type Form = form.Form

type Sort struct {
	Form Form
	Attr SortAttr
}

type SortAttr interface {
	Level() int
	Parent() SortAttr
	LessEqual(dst SortAttr) bool
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
