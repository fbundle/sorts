package el2

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
)

type runtimeParser struct {
	parseName   func(Name) Sort
	listParsers ordered_map.OrderedMap[Name, ListParseFunc]
}

func (u runtimeParser) parse(node Form) AlmostSort {
	switch node := node.(type) {
	case Name:
		return ActualSort{sort: u.parseName(node)}
	case List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := u.listParsers.Get(head); ok {
			return listParser(u.parse, node)
		} else { // by default, parse as beta reduction (function call)
			return ListParseBeta(u.parse, node)
		}
	default:
		panic("parse error")
	}
}

func (u runtimeParser) newListParser(head Name, parseList ListParseFunc) runtimeParser {
	if _, ok := u.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	u.listParsers = u.listParsers.Set(head, parseList)
	return u
}
