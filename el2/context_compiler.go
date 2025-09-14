package el2

import (
	"encoding/json"
	"log"

	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/form"
)

func logCompile(node form.Form, sort el_sorts.Sort) {
	toString := func(o any) string {
		b, err := json.Marshal(o)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	log.Printf("compiled %s -> %s", toString(node), toString(sort))
}

func (ctx Context) Compile(node form.Form) el_sorts.Sort {
	switch node := node.(type) {
	case form.Name:
		sort := ctx.Get(node)
		logCompile(node, sort)
		return sort
	case form.List:
		if len(node) == 0 {
			panic("empty list")
		}
		if head, ok := node[0].(form.Name); ok {
			if listParser, ok := ctx.listCompiler.Get(head); ok {
				sort := listParser(ctx, node)
				logCompile(node, sort)
				return sort
			}
		}
		// use default
		sort := ctx.defaultListCompiler(ctx, node)
		logCompile(node, sort)
		return sort
	default:
		panic("parse_error")
	}
}

func (ctx Context) WithListCompiler(name form.Name, compileFunc func(form.Name) el_sorts.ListCompileFunc) Context {
	ctx.listCompiler = ctx.listCompiler.Set(name, compileFunc(name))
	return ctx
}
