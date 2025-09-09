package sorts

import (
	"strconv"

	"github.com/fbundle/sorts/expr"
)

func NewAtom(level int, term expr.Term, parent Sort) Atom {
	if parent != nil && Level(parent) != level+1 {
		panic("type_error make")
	}
	return Atom{
		level:  level,
		term:   term,
		parent: parent,
	}
}
func defaultSort(sort Sort, level int) Sort {
	if sort != nil {
		if Level(sort) != level {
			panic("type_error level")
		}
		return sort
	}
	return Atom{
		level:  level,
		term:   expr.Term(defaultTerm + "_" + strconv.Itoa(level)),
		parent: nil,
	}
}

// dummyTerm - make a dummy term of type parent
func dummyTerm(parent Sort, term expr.Term) Sort {
	return NewAtom(Level(parent)-1, term, parent)
}

type Atom struct {
	level  int
	term   expr.Term
	parent Sort
}

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		repr:   s.term,
		level:  s.level,
		parent: defaultSort(s.parent, s.level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				if s.term == InitialTerm || d.term == TerminalTerm {
					return true
				}
				if s.term == d.term {
					return true
				}
				_, ok := lessEqualMap[rule{s.term, d.term}]
				return ok
			default:
				return false
			}
		},
	}
}

type rule struct {
	src expr.Term
	dst expr.Term
}

var lessEqualMap = make(map[rule]struct{})

func AddRule(src expr.Term, dst expr.Term) {
	lessEqualMap[rule{src, dst}] = struct{}{}
}
