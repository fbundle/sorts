package sorts

import "github.com/fbundle/sorts/form"

const (
	InitialName  form.Name = "Unit" // initial object
	TerminalName form.Name = "Any"  // terminal object
	DefaultName  form.Name = "Type"
)

func defaultSort(level int) Sort {
	return newAtom(level, DefaultName, func(level int) form.Name {
		return DefaultName
	})
}

func newAtom(level int, repr form.Name, ancestor func(int) form.Name) Sort {
	parentLevel := level + 1
	parentRepr := ancestor(parentLevel)
	return Atom{
		level: level,
		name:  repr,
		parent: func() Sort {
			return newAtom(parentLevel, parentRepr, ancestor)
		},
	}
}

func newTerm(termRepr form.Name, parent Sort) Sort {
	return Atom{
		level:  Level(parent) - 1,
		name:   termRepr,
		parent: func() Sort { return parent },
	}
}

type Atom struct {
	level  int
	name   form.Name
	parent func() Sort
}

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		repr:   s.name,
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				if s.name == InitialName || d.name == TerminalName {
					return true
				}
				if s.name == d.name {
					return true
				}

				_, ok := lessEqualMap[rule{s.name, d.name}]
				return ok

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
