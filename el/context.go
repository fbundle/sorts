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
type Sort = sorts.Sort
type ListParser = sorts.ListParser
type Context = sorts.Context

const (
	InitialName  = "Unit"
	TerminalName = "Any"
	DefaultName  = "Type"
)

type EL struct {
	frame       ordered_map.OrderedMap[Name, Sort]
	listParsers ordered_map.OrderedMap[Name, ListParser]
}

func (ctx EL) Get(name Name) Sort {
	if sort, ok := ctx.frame.Get(name); ok {
		return sort
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (ctx EL) Set(name Name, sort Sort) Context {
	ctx.frame = ctx.frame.Set(name, sort)
	return ctx
}

func (ctx EL) Parse(node Form) (Context, Sort) {
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
	}

	//TODO implement me
	panic("implement me")
}

func (ctx EL) AddListParser(listParser sorts.ListParser) sorts.Context {
	//TODO implement me
	panic("implement me")
}

func (ctx EL) LessEqual(src Form, dst Form) bool {
	//TODO implement me
	panic("implement me")
}

var _ sorts.Context = EL{}
