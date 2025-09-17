package el

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func (ctx Context) Get(name form.Name) sorts.Sort1 {
	if s, ok := ctx.frame.Get(name); ok {
		return s
	}
	if s, ok := ctx.universe.GetBuiltin(name); ok {
		return s
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (ctx Context) Set(name form.Name, sort sorts.Sort1) sorts.Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}
func (ctx Context) Del(name form.Name) sorts.Context {
	ctx.frame = ctx.frame.Del(name)
	return ctx
}
