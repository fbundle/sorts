package el

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Name = form.Name
type List = form.List
type Form = form.Form

const (
	InitialName  = "Unit"
	TerminalName = "Any"
	DefaultName  = "Type"
)

type rule struct {
	src Name
	dst Name
}

func (r rule) Cmp(s rule) int {
	if c := cmp.Compare(r.src, s.src); c != 0 {
		return c
	}
	return cmp.Compare(r.dst, s.dst)
}

type Context struct {
	frame       ordered_map.OrderedMap[Name, sorts.Sort]
	listParsers ordered_map.OrderedMap[Name, sorts.ListParseFunc]
	nameRule    ordered_map.Map[rule]
}

func (ctx Context) Get(name Name) sorts.Sort {
	if sort, ok := ctx.frame.Get(name); ok {
		return sort
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (ctx Context) Set(name Name, sort sorts.Sort) sorts.Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}

func (ctx Context) Parse(node Form) (sorts.Context, sorts.Sort) {
	switch node := node.(type) {
	case Name:
		if sort, ok := ctx.frame.Get(node); ok {
			return ctx, sort
		}
		// parse builtin
		for _, name := range []Name{InitialName, TerminalName, DefaultName} {
			prefix := name + "_"
			if strings.HasPrefix(string(node), string(prefix)) {
				levelStr := strings.TrimPrefix(string(node), string(prefix))
				if level, err := strconv.Atoi(levelStr); err == nil {
					parent := sorts.NewChain(DefaultName, level+1)
					return ctx, sorts.NewTerm(
						node,
						func(ctx sorts.Context) sorts.Sort {
							return parent
						},
					)
				}
			}
		}
		panic(fmt.Errorf("name_not_found: %s", node))
	case List:
		if len(node) == 0 {
			panic(fmt.Errorf("list_empty"))
		}
		if head, ok := node[0].(Name); ok {
			if listParse, ok := ctx.listParsers.Get(head); ok {
				return listParse(ctx, node[1:])
			}
		}
		panic(fmt.Errorf("parse_error: %v", node))
	default:
		panic(fmt.Errorf("parse_error: %v", node))
	}
}

func (ctx Context) AddListParser(listParser sorts.ListParser) sorts.Context {
	ctx.listParsers = ctx.listParsers.Set(listParser.Command, listParser.ListParse)
	return ctx
}

func (ctx Context) LessEqual(src Form, dst Form) bool {
	s, ok1 := src.(Name)
	d, ok2 := dst.(Name)
	if ok1 && s == InitialName {
		return true
	}
	if ok2 && d == TerminalName {
		return true
	}

}

var _ sorts.Context = Context{}
