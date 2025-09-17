package sorts5

import "github.com/fbundle/sorts/form"

type Name = form.Name
type List = form.List
type Form = form.Form

type Frame interface {
	Get(name string) Sort
	Set(name string, sort Sort) Frame
	Del(name string) Frame
}

type Sort interface {
	Compile(frame Frame) Sort
	Form() Form
	Level(frame Frame) int
	Parent(frame Frame) Sort
	LessEqual(frame Frame, d Sort) bool
}
