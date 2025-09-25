package sorts_context

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/persistent/seq"
	"github.com/fbundle/sorts/sorts"
)

type Form = form.Form
type Name = form.Name
type List = form.List
type Sort = sorts.Sort
type Code = sorts.Code
type Context struct {
	builtin seq.Seq[func(Name) (sorts.Sort, bool)]
	frame   ordered_map.OrderedMap[Name, Sort]
	Univ
}

func (c Context) WithBuiltin(get func(Name) (sorts.Sort, bool)) Context {
	c.builtin = c.builtin.PushBack(get)
	return c
}

func (c Context) Get(name sorts.Name) sorts.Sort {
	if value, ok := c.frame.Get(name); ok {
		return value
	}
	for _, builtin := range c.builtin.Iter {
		if value, ok := builtin(name); ok {
			return value
		}
	}
	if value, ok := c.Univ.builtinNameGet(name); ok {
		return value
	}
	panic(fmt.Errorf("name_not_found: %s", name))
}

func (c Context) Set(name sorts.Name, sort sorts.Sort) sorts.Context {
	c.frame = c.frame.Set(name, sort)
	return c
}
