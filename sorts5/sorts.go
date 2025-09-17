package sorts5

import (
	"github.com/fbundle/sorts/form"
)

type Form = form.Form

type Sort struct {
	Form      Form
	Level     func() int
	Parent    func() Sort
	LessEqual func(dst Sort) bool
}

type Mode string

const (
	ModeComp  Mode = "COMP"  // type checking
	ModeEval  Mode = "EVAL"  // type checking and evaluation
	ModeDebug Mode = "DEBUG" // type checking and print everything
)

// Compiler - recursive compilation
type Compiler interface {
	Compile(form Form) Sort
}

// Frame - name binding
type Frame interface {
	Get(name string) Sort
	Set(name string, sort Sort) Context
	Del(name string) Context
}
type Context struct {
}
