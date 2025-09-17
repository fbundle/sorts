package sorts5

import (
	"strconv"
	"strings"
)

const (
	initialName  = "Unit"
	terminalName = "Any"
)

func parseBuiltin(name Name) (Sort, bool) {
	nameStr := string(name)
	// parse builtin name such as Unit and Any
	for _, builtinName := range []string{initialName, terminalName} {
		prefix := builtinName + "_"
		if strings.HasPrefix(nameStr, prefix) {
			levelStr := strings.TrimPrefix(nameStr, prefix)
			if level, err := strconv.Atoi(levelStr); err == nil {
				return newChain(Name(builtinName), level), true
			}
		}
	}
	return nil, false
}

func newChain(name Name, level int) Atom {
	return Atom{
		form: name,
		level: func() int {
			return level
		},
		parent: func() Sort {
			return newChain(name, level+1)
		},
	}
}

type Atom struct {
	form   Form
	level  func() int
	parent func() Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Level() int {
	return s.level()
}

func (s Atom) Parent() Sort {
	return s.parent()
}

func (s Atom) LessEqual(dst Sort) bool {
	// compare form
	sName, ok1 := isName(s.Form())
	dName, ok2 := isName(dst.Form())
	if ok1 && sName == initialName {
		return true
	}
	if ok2 && dName == terminalName {
		return true
	}
	if ok1 && ok2 {
		if _, ok3 := ruleMap[[2]Name{sName, dName}]; ok3 {
			return true
		}
	}
	return false
}

var ruleMap = map[[2]Name]struct{}{}

func AddRule(src Name, dst Name) {
	ruleMap[[2]Name{src, dst}] = struct{}{}
}
