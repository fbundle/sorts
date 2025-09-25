package sorts

import "github.com/fbundle/sorts/persistent/ordered_map"

type ListParseFunc = func(parse func(form Form) Code, list List) Code
type NameParseFunc = func(name Name) Code

type parser struct {
	builtinNameParseFunc func(name Name) (Code, bool)

	finalListParseFunc ListParseFunc
	finalNameParseFunc NameParseFunc
	listParseFuncMap   ordered_map.OrderedMap[Name, ListParseFunc]
}

func (p parser) Parse(form Form) Code {
	switch f := form.(type) {
	case Name:
		if p.builtinNameParseFunc != nil {
			if code, ok := p.builtinNameParseFunc(f); ok {
				return code
			}
		}
		return p.finalNameParseFunc(f)
	case List:
		if len(f) == 0 {
			panic("parse_error")
		}
		if cmd, ok := f[0].(Name); ok {
			if parseFunc, ok := p.listParseFuncMap.Get(cmd); ok {
				return parseFunc(p.Parse, f[1:])
			}
		}
		return p.finalListParseFunc(p.Parse, f)
	}
	panic("parse_error")
}

var p = parser{}

func setFinalListParseFunc(finalListParseFunc ListParseFunc) {
	p.finalListParseFunc = finalListParseFunc
}
func setFinalNameParseFunc(finalNameParseFunc NameParseFunc) {
	p.finalNameParseFunc = finalNameParseFunc
}
func addListParseFunc(cmd Name, listParseFunc ListParseFunc) {
	p.listParseFuncMap = p.listParseFuncMap.Set(
		cmd,
		listParseFunc,
	)
}

func MakeParser(builtinNameParseFunc func(name Name) (Code, bool)) func(form Form) Code {
	return parser{
		builtinNameParseFunc: builtinNameParseFunc,
		finalListParseFunc:   p.finalListParseFunc,
		finalNameParseFunc:   p.finalNameParseFunc,
		listParseFuncMap:     p.listParseFuncMap,
	}.Parse
}
