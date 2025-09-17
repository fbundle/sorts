package el

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func (ctx Context) Form(s any) sorts.Form {
	return ctx.universe.Form(s)
}

func (ctx Context) Level(s sorts.Sort) int {
	return ctx.universe.Level(s)
}

func (ctx Context) Parent(s sorts.Sort) sorts.Sort {
	return ctx.universe.Parent(s)
}

func (ctx Context) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return ctx.universe.LessEqual(x, y)
}

func (ctx Context) LessEqualBasic(x sorts.Sort, y sorts.Sort) bool {
	return ctx.universe.LessEqualBasic(x, y)
}

func (ctx Context) NewTerm(name form.Form, parent sorts.Sort) sorts.Atom {
	return ctx.universe.NewTerm(name, parent)
}
func (ctx Context) Initial(level int) sorts.Sort {
	return ctx.universe.Initial(level)
}

func (ctx Context) Terminal(level int) sorts.Sort {
	return ctx.universe.Terminal(level)
}
