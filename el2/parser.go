package el2

import (
	"github.com/fbundle/sorts/el2/el_almost_sort"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Parser struct {
	parseName   func(Name) sorts.Sort
	listParsers ordered_map.OrderedMap[Name, ListParseFunc]
}

func (p Parser) Parse(node Form) el_almost_sort.AlmostSort {
	switch node := node.(type) {
	case Name:
		return el_almost_sort.ActualSort{Sort: p.parseName(node)}
	case List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := p.listParsers.Get(head); ok {
			return listParser(p.Parse, node)
		} else { // by default, Parse as beta reduction (function call)
			return el_almost_sort.ListParseBeta(p.Parse, node)
		}
	default:
		panic("Parse error")
	}
}

func (p Parser) newListParser(head Name, parseList ListParseFunc) Parser {
	if _, ok := p.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	p.listParsers = p.listParsers.Set(head, parseList)
	return p
}
