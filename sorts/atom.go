package sorts

const (
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
)

func MakeAtom(level int, name string, parentName string) Sort {
	return Atom{
		level:   level,
		name:    name,
		_parent: _atomParent{name: parentName, sort: nil},
	}
}
func MakeTerm(termName string, parent Sort) Sort {
	return Atom{
		level:   Level(parent) - 1,
		name:    termName,
		_parent: _atomParent{name: "", sort: parent},
	}
}

type Atom struct {
	level   int
	name    string
	_parent _atomParent
}

func (s Atom) parent() Sort {
	if s._parent.sort != nil {
		return s._parent.sort
	}
	parentLevel := s.level + 1
	parentName := s._parent.name
	grandParentName := s._parent.name
	return MakeAtom(parentLevel, parentName, grandParentName)
}

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		name:   s.name,
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

// _atomParent - parent of atom - must be either a name or a sort
type _atomParent struct {
	name string
	sort Sort
}

type rule struct {
	src string
	dst string
}

var lessEqualMap = make(map[rule]struct{})

func AddRule(src string, dst string) {
	lessEqualMap[rule{src, dst}] = struct{}{}
}
