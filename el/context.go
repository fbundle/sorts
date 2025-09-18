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

type Context struct {
	frame       ordered_map.OrderedMap[Name, sorts.Sort]
	listParsers ordered_map.OrderedMap[Name, sorts.ListCompileFunc]
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

func (ctx Context) Parse(node Form) sorts.Sort {
	switch node := node.(type) {
	case Name:
		// all names should be either builtin or linked to a Sort
		if sort, ok := ctx.frame.Get(node); ok {
			return sort
		}
		// parse builtin
		for _, name := range []Name{InitialName, TerminalName, DefaultName} {
			prefix := name + "_"
			if strings.HasPrefix(string(node), string(prefix)) {
				levelStr := strings.TrimPrefix(string(node), string(prefix))
				if level, err := strconv.Atoi(levelStr); err == nil {
					parent := sorts.NewChain(DefaultName, level+1)
					return sorts.NewTerm(node, parent)
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
		return sorts.DefaultCompileFunc(ctx, node)
	default:
		panic(fmt.Errorf("parse_error: %v", node))
	}
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
	if ok1 && ok2 {
		if _, ok := ctx.nameRule.Get(rule{s, d}); ok {
			return true
		}
	}
	return false
}

var _ sorts.Context = Context{}

func (ctx Context) Init() Context {
	for cmd, listParseFunc := range sorts.ListCompileFuncMap {
		ctx = ctx.AddListParseFunc(cmd, listParseFunc)
	}
	return ctx
}
func (ctx Context) AddListParseFunc(cmd Name, listParse sorts.ListCompileFunc) Context {
	ctx.listParsers = ctx.listParsers.Set(cmd, listParse)
	return ctx
}

func (ctx Context) AddLessEqualRule(src Name, dst Name) Context {
	ctx.nameRule = ctx.nameRule.Set(rule{src, dst})
	return ctx
}

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
