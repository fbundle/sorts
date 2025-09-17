package el

import (
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

type Context struct {
	frame       ordered_map.OrderedMap[Name, sorts.Sort]
	listParsers ordered_map.OrderedMap[Name, sorts.ListParser]
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
				return listParse.ListParse(ctx, node[1:])
			}
		}
		panic(fmt.Errorf("parse_error: %v", node))
	default:
		panic(fmt.Errorf("parse_error: %v", node))
	}
}

func (ctx Context) AddListParser(listParser sorts.ListParser) sorts.Context {
	//TODO implement me
	panic("implement me")
}

func (ctx Context) LessEqual(src Form, dst Form) bool {
	//TODO implement me
	panic("implement me")
}

var _ sorts.Context = Context{}
