package sorts_parser

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Form = form.Form
type Name = form.Name
type List = form.List
type Sort = sorts.Sort
type Code = sorts.Code

type ListParseFunc = func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code
type NameParseFunc = func(name sorts.Name) sorts.Code

type Parser struct {
	finalListParseFunc ListParseFunc
	finalNameParseFunc NameParseFunc
	listParseFuncMap   ordered_map.OrderedMap[sorts.Name, ListParseFunc]
}

func (p Parser) Parse(form sorts.Form) sorts.Code {
	switch f := form.(type) {
	case sorts.Name:
		return p.finalNameParseFunc(f)
	case sorts.List:
		if len(f) == 0 {
			panic("parse_error")
		}
		if cmd, ok := f[0].(sorts.Name); ok {
			if parseFunc, ok := p.listParseFuncMap.Get(cmd); ok {
				return parseFunc(p.Parse, f[1:])
			}
		}
		return p.finalListParseFunc(p.Parse, f)
	}
	panic("parse_error")
}
