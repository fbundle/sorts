package parser

import (
	"github.com/fbundle/sorts/el2"
	"github.com/fbundle/sorts/persistent/ordered_map"
)

type parser struct {
	parseName   func(el2.Name) el2.Sort
	listParsers ordered_map.OrderedMap[el2.Name, el2.ListParseFunc]
}

func (u parser) parse(node el2.Form) el2.AlmostSort {
	switch node := node.(type) {
	case el2.Name:
		return el2.ActualSort{sort: u.parseName(node)}
	case el2.List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(el2.Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := u.listParsers.Get(head); ok {
			return listParser(u.parse, node)
		} else { // by default, parse as beta reduction (function call)
			return el2.ListParseBeta(u.parse, node)
		}
	default:
		panic("parse error")
	}
}

func (u parser) newListParser(head el2.Name, parseList el2.ListParseFunc) parser {
	if _, ok := u.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	u.listParsers = u.listParsers.Set(head, parseList)
	return u
}
