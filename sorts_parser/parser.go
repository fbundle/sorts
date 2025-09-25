package sorts_parser

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type ListParseFunc = func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code
type NameParseFunc = func(name sorts.Name) sorts.Code

type parser struct {
	builtinNameParseFunc func(name sorts.Name) (sorts.Code, bool)

	finalListParseFunc ListParseFunc
	finalNameParseFunc NameParseFunc
	listParseFuncMap   ordered_map.OrderedMap[sorts.Name, ListParseFunc]
}

func (p parser) Parse(form sorts.Form) sorts.Code {
	switch f := form.(type) {
	case sorts.Name:
		if p.builtinNameParseFunc != nil {
			if code, ok := p.builtinNameParseFunc(f); ok {
				return code
			}
		}
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

var p = parser{}

func SetFinalListParseFunc(finalListParseFunc ListParseFunc) {
	if p.finalListParseFunc != nil {
		panic("set_twice")
	}
	p.finalListParseFunc = finalListParseFunc
}
func SetFinalNameParseFunc(finalNameParseFunc NameParseFunc) {
	if p.finalNameParseFunc != nil {
		panic("set_twice")
	}
	p.finalNameParseFunc = finalNameParseFunc
}
func AddListParseFunc(cmd sorts.Name, listParseFunc ListParseFunc) {
	p.listParseFuncMap = p.listParseFuncMap.Set(
		cmd,
		listParseFunc,
	)
}

func MakeParser(builtinNameParseFunc func(name sorts.Name) (sorts.Code, bool)) func(form sorts.Form) sorts.Code {
	return parser{
		builtinNameParseFunc: builtinNameParseFunc,
		finalListParseFunc:   p.finalListParseFunc,
		finalNameParseFunc:   p.finalNameParseFunc,
		listParseFuncMap:     p.listParseFuncMap,
	}.Parse
}
