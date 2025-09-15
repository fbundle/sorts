package el

import (
	"fmt"

	"github.com/fbundle/sorts/el/el_sorts"
	"github.com/fbundle/sorts/form"
)

func (ctx Context) Get(name form.Name) el_sorts.Sort {
	if s, ok := ctx.frame.Get(name); ok {
		return s
	}
	if s, ok := ctx.universe.GetBuiltin(name); ok {
		return s
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (ctx Context) Set(name form.Name, sort el_sorts.Sort) el_sorts.Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}
func (ctx Context) Del(name form.Name) el_sorts.Context {
	ctx.frame = ctx.frame.Del(name)
	return ctx
}
