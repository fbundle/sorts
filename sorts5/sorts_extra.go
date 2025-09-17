package sorts5

import (
	"strconv"
	"strings"
)

const (
	InitialName  Name = "Unit"
	TerminalName Name = "Any"
	DefaultName  Name = "Type"
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

type ListParseFunc = func(List) Sort

var listParseFuncMap = map[Name]ListParseFunc{}

func AddListParseFunc(name Name, fn ListParseFunc) {
	listParseFuncMap[name] = fn
}

func Parse(frame Frame, form Form) Sort {
	switch form := form.(type) {
	case Name:
		for _, builtinName := range []Name{InitialName, TerminalName} {
			prefix := string(builtinName + "_")
			if strings.HasPrefix(string(form), prefix) {
				levelStr := string(form[len(prefix):])
				if level, err := strconv.Atoi(levelStr); err == nil {
					return Atom{
						form: form,
						level: func(ctx Frame) int {
							return level
						},
						parent: func(ctx Frame) Sort {
							return NewChain(DefaultName, level+1)
						},
					}
				}

			}
		}

	case List:
	}

}
