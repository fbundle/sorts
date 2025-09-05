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
	parent := s.parent
	if parent == nil {
		// default parent
		parent = Atom{
			level:  s.level,
			name:   defaultName + "_" + strconv.Itoa(s.level),
			parent: nil,
		}
	}

	lessEqual := func(dst Sort) bool {
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

		default:
			panic("type_error - should catch all types")
		}
	}

	return SortView{
		Level:     s.level,
		Name:      s.name,
		Parent:    parent,
		LessEqual: lessEqual,
	}
}

type rule struct {
	src string
	dst string
}

var lessEqualMap = make(map[rule]struct{})
