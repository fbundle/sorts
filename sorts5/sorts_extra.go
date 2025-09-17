package sorts5

import "github.com/fbundle/sorts/form"

const (
	InitialName  Name = "Unit"
	TerminalName Name = "Any"
)

var ruleMap = map[[2]Name]struct{}{}

func AddRule(src Name, dst Name) {
	ruleMap[[2]Name{src, dst}] = struct{}{}
}

func FallbackLessEqual(src Form, dst Form) bool {
	s, ok1 := src.(Name)
	d, ok2 := dst.(Name)
	if ok1 && s == InitialName {
		return true
	}
	if ok2 && d == TerminalName {
		return true
	}
	if ok1 && ok2 {
		if _, ok3 := ruleMap[[2]Name{s, d}]; ok3 {
			return true
		}
	}
	return false
}

type ListParseFunc = func(form.List) Sort

var listParseFuncMap = map[Name]ListParseFunc{}

func AddListParseFunc(name Name, fn ListParseFunc) {
	listParseFuncMap[name] = fn
}
