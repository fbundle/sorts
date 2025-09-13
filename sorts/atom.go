package sorts

func NewAtomChain(level int, chainName func(int) Name) Atom {
	return Atom{
		level: level,
		name:  chainName(level),
		parent: func() Sort {
			return NewAtomChain(level+1, chainName)
		},
	}
}
func NewAtomTerm(a SortAttr, name Name, parent Sort) Atom {
	return Atom{
		level: a.Level(parent) - 1,
		name:  name,
		parent: func() Sort {
			return parent
		},
	}
}

type Atom struct {
	level  int
	name   Name
	parent func() Sort
}

func (s Atom) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   s.name,
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				return a.NameLessEqual(s.name, d.name)
			default:
				return false
			}
		},
	}
}
