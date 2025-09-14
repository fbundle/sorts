package el2

import (
	"github.com/fbundle/sorts/el2/frame"
	"github.com/fbundle/sorts/el2/sort_universe"
	"github.com/fbundle/sorts/persistent/ordered_map"
)

type Runtime struct {
	frame        el_frame.Frame
	sortUniverse el_sort_universe.SortUniverse

	listParsers ordered_map.OrderedMap[Name, ListParseFunc]
}

func (r Runtime) Parse(node Form) AlmostSort {
	switch node := node.(type) {
	case Name:
		if sort, ok := r.frame.Get(node); ok {
			return ActualSort{sort}
		}
		if sort, ok := r.sortUniverse.ParseBuiltin(node); ok {
			return ActualSort{sort}
		}
		panic("name_not_found")
	case List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := r.listParsers.Get(head); ok {
			return listParser(r.Parse, node)
		} else { // by default, Parse as beta reduction (function call)
			return ListParseBeta(r.Parse, node)
		}
	default:
		panic("Parse error")
	}
}

func (r Runtime) newListParser(head Name, parseList ListParseFunc) Runtime {
	if _, ok := r.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	r.listParsers = r.listParsers.Set(head, parseList)
	return r
}
