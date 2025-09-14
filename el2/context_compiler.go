package el2

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/form"
)

func (ctx Context) ToString(o any) string {
	switch v := o.(type) {
	case string:
		return v
	case form.Form:
		return strings.Join(v.Marshal("(", ")"), " ")
	case el_sorts.Sort:
		f := ctx.ToString(ctx.Form(v))
		t := strings.Join(ctx.Form(ctx.Parent(v)).Marshal("(", ")"), " ")
		l := ctx.Level(v)
		return fmt.Sprintf("(form %s - type %s - level %d)", f, t, l)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func logCompile(ctx Context, node form.Form, sort el_sorts.Sort) {
	log.Printf("compiled %s from %s\n", ctx.ToString(sort), ctx.ToString(node))
}

func (ctx Context) Compile(node form.Form) el_sorts.Sort {
	switch node := node.(type) {
	case form.Name:
		sort := ctx.Get(node)
		logCompile(ctx, node, sort)
		return sort
	case form.List:
		if len(node) == 0 {
			panic("empty list")
		}
		if head, ok := node[0].(form.Name); ok {
			if listParser, ok := ctx.listCompiler.Get(head); ok {
				sort := listParser(ctx, node)
				logCompile(ctx, node, sort)
				return sort
			}
		}
		// use default
		sort := ctx.defaultListCompiler(ctx, node)
		logCompile(ctx, node, sort)
		return sort
	default:
		panic("parse_error")
	}
}

func (ctx Context) WithListCompiler(name form.Name, compileFunc func(form.Name) el_sorts.ListCompileFunc) Context {
	ctx.listCompiler = ctx.listCompiler.Set(name, compileFunc(name))
	return ctx
}
