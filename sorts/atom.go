package sorts

func newAtomChain(level int, chainName func(int) Name) Atom {
	return Atom{
		level: level,
		name:  chainName(level),
		parent: func() Sort {
			return newAtomChain(level+1, chainName)
		},
	}
}
func newAtomTerm(u Universe, name Name, parent Sort) Atom {
	return Atom{
		level: u.Level(parent) - 1,
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

func (s Atom) sortAttr() sortAttr {
	return sortAttr{
		repr:   s.name,
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(u Universe, dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				return u.lessEqual(s.name, d.name)
			default:
				return false
			}
		},
	}
}
