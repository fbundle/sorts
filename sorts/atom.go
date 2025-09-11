package sorts

import (
	"fmt"
	"strconv"
)

const (
	nameWithType = false
)

func NewAtom(level int, name string, parent Sort) Atom {
	if parent != nil && Level(parent) != level+1 {
		panic("type_error make")
	}
	return Atom{
		level:  level,
		name:   name,
		parent: parent,
	}
}

// dummyTerm - make a dummy term of type parent
func dummyTerm(parent Sort, name string) Sort {
	newName := name
	if nameWithType {
		newName = fmt.Sprintf("(%s: %s)", name, Name(parent))
	}
	return NewAtom(Level(parent)-1, newName, parent)
}

type Atom struct {
	level  int
	name   string
	parent Sort
}

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		name:   s.name,
		level:  s.level,
		parent: defaultSort(s.parent, s.level+1),
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

func defaultSort(sort Sort, level int) Sort {
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
