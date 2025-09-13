package sorts

func NewAtomChain(level int, chainName func(int) Name) Sort {
	return Atom{
		level: level,
		form:  chainName(level),
		parent: func() Sort {
			return NewAtomChain(level+1, chainName)
		},
	}
}
func NewAtomTerm(a SortAttr, form Form, parent Sort) Sort {
	return Atom{
		level: a.Level(parent) - 1,
		form:  form,
		parent: func() Sort {
			return parent
		},
	}
}

type Atom struct {
	level  int
	form   Form
	parent func() Sort
}

func (s Atom) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   s.form,
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				if s.level != d.level {
					return false
				}
				sname, ok1 := s.form.(Name)
				dname, ok2 := d.form.(Name)
				if ok1 && ok2 {
					return a.NameLessEqual(sname, dname)
				}
				return false
			default:
				return false
			}
		},
	}
}
