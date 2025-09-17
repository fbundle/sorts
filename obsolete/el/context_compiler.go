package el

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func (ctx Context) ToString(o any) string {
	switch v := o.(type) {
	case string:
		return v
	case form.Form:
		return strings.Join(v.Marshal("(", ")"), " ")
	case sorts.Sort1:
		f := ctx.ToString(ctx.Form(v))
		t := strings.Join(ctx.Form(ctx.Parent(v)).Marshal("(", ")"), " ")
		l := ctx.Level(v)

		if ctx.Mode() == sorts.ModeComp {
			return fmt.Sprintf("(type %s - level %d)", t, l)
		} else {
			return fmt.Sprintf("(form %s - type %s - level %d)", f, t, l)
		}
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func logCompile(ctx Context, node form.Form, sort sorts.Sort1) {
	if ctx.Mode() == sorts.ModeDebug {
		log.Printf("DEBUG: compiled %s from %s\n", ctx.ToString(sort), ctx.ToString(node))
	}
}

func (ctx Context) Compile(node form.Form) sorts.Sort1 {
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

func (ctx Context) WithListCompiler(name form.Name, compileFunc func(form.Name) sorts.ListCompileFunc) Context {
	// TODO - instead of assuming list starts with name
	// how about we have

	ctx.listCompiler = ctx.listCompiler.Set(name, compileFunc(name))
	return ctx
}
