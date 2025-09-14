package el2_parser

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Parser struct {
	ParseName   func(form.Name) sorts.Sort
	listParsers ordered_map.OrderedMap[form.Name, el2_almost_sort.ListParseFunc]
}

func (p Parser) Parse(node form.Form) el2_almost_sort.AlmostSort {
	switch node := node.(type) {
	case form.Name:
		return el2_almost_sort.ActualSort{Sort: p.ParseName(node)}
	case form.List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(form.Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := p.listParsers.Get(head); ok {
			return listParser(p.Parse, node)
		} else { // by default, Parse as beta reduction (function call)
			return el2_almost_sort.ListParseBeta(p.Parse, node)
		}
	default:
		panic("Parse error")
	}
}

func (p Parser) WithListParser(head form.Name, parseList el2_almost_sort.ListParseFuncWithHead) Parser {
	if _, ok := p.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	p.listParsers = p.listParsers.Set(head, parseList(head))
	return p
}
