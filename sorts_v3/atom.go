package sorts

import "strconv"

func Unit(level int) Atom {
	return NewAtom(
		level,
		InitialName,
		nil,
	)
}
func Any(level int) Atom {
	return NewAtom(
		level,
		TerminalName,
		nil,
	)
}

func NewAtom(level int, name string, parent WithSort) Atom {
	if parent != nil && Level(parent) != level+1 {
		panic("type_error make")
	}
	return Atom{
		level:  level,
		name:   name,
		parent: parent,
	}
}

type Atom struct {
	level  int
	name   string
	parent WithSort
}

func (s Atom) attr() sortAttr {
	return sortAttr{
		level:  s.level,
		name:   s.name,
		parent: defaultSort(s.parent, s.level+1),
		lessEqual: func(dst WithSort) bool {
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

func defaultSort(sort WithSort, level int) WithSort {
	if sort != nil {
		if Level(sort) != level {
			panic("type_error Level")
		}
		return sort
	}
	return Atom{
		level:  level,
		name:   defaultName + "_" + strconv.Itoa(level),
		parent: nil,
	}
}

type rule struct {
	src string
	dst string
}

var lessEqualMap = make(map[rule]struct{})

func AddRule(src string, dst string) {
	lessEqualMap[rule{src, dst}] = struct{}{}
}
