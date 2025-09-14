package el2

import (
	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/form"
)

func (ctx Context) Compile(node form.Form) (almost_sort_extra.Context, almost_sort_extra.Sort) {
	switch node := node.(type) {
	case form.Name:
		return ctx, ctx.Get(node)
	case form.List:
		if len(node) == 0 {
			panic("empty list")
		}
		if head, ok := node[0].(form.Name); ok {
			if listParser, ok := ctx.listCompiler.Get(head); ok {
				return listParser(ctx, node)
			}
		}
		// use default
		return ctx.defaultListCompiler(ctx, node)
	default:
		panic("parse_error")
	}
}

func (ctx Context) WithListCompiler(name form.Name, compileFunc func(form.Name) almost_sort_extra.ListCompileFunc) Context {
	ctx.listCompiler = ctx.listCompiler.Set(name, compileFunc(name))
	return ctx
}
