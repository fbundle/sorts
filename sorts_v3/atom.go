package sorts

import "strconv"

const (
	defaultName  = "type"
	initialName  = "unit" // initial object
	terminalName = "any"  // terminal object
)

var _ Sort = Atom{}

func NewAtom(level int, name string, parent Sort) Atom {
	if parent != nil && parent.View().Level != level+1 {
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
	parent Sort
}

func (s Atom) View() SortView {
	return SortView{
		Level: s.level,
		Name:  s.name,
		Parent: InhabitedSort{
			Sort:  defaultSort(s.parent, s.level+1),
			Child: s,
		},
		LessEqual: func(dst Sort) bool {
			switch d := dst.(type) {

			case Atom:
				if s.level != d.level {
					return false
				}
				if s.name == initialName || d.name == terminalName {
					return true
				}
				if s.name == d.name {
					return true
				}
				_, ok := lessEqualMap[rule{s.name, d.name}]
				return ok
			case Arrow:
				return false
			default:
				panic("type_error - should catch all types")
			}
		},
	}
}

func defaultSort(sort Sort, level int) Sort {
	if sort != nil {
		if sort.View().Level != level {
			panic("type_error level")
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
