package sorts

import "github.com/fbundle/sorts/form"

const (
	InitialName  form.Name = "Unit" // initial object
	TerminalName form.Name = "Any"  // terminal object
	DefaultName  form.Name = "Type"
)

func defaultSort(level int) Sort {
	return NewAtom(level, DefaultName, func(level int) form.Form {
		return DefaultName
	})
}

func NewAtom(level int, repr form.Form, ancestor func(int) form.Form) Sort {
	parentLevel := level + 1
	parentRepr := ancestor(parentLevel)
	return Atom{
		level: level,
		repr:  repr,
		parent: func() Sort {
			return NewAtom(parentLevel, parentRepr, ancestor)
		},
	}
}

func NewTerm(termRepr form.Form, parent Sort) Sort {
	return Atom{
		level:  Level(parent) - 1,
		repr:   termRepr,
		parent: func() Sort { return parent },
	}
}

type Atom struct {
	level  int
	repr   form.Form
	parent func() Sort
}

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		repr:   s.repr,
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				if s.repr == InitialName || d.repr == TerminalName {
					return true
				}
				if s.repr == d.repr {
					return true
				}

				srepr, ok1 := s.repr.(form.Name)
				drepr, ok2 := d.repr.(form.Name)
				if ok1 && ok2 {
					_, ok := lessEqualMap[rule{srepr, drepr}]
					return ok
				}
				return false

			default:
				return false
			}
		},
	}
}

type rule struct {
	src form.Name
	dst form.Name
}

var lessEqualMap = make(map[rule]struct{})

func AddRule(src form.Name, dst form.Name) {
	lessEqualMap[rule{src, dst}] = struct{}{}
}
