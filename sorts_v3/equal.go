package sorts

import "fmt"

type Equal struct {
	A Sort
	B Sort
}

func (s Equal) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("Eq(%s, %s)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Equal:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

func (s Equal) Refl(x Sort) Sort {
	// reflexive
	MustTermOf(x, s.A)
	MustTermOf(x, s.B)
	return dummyTerm(s, fmt.Sprintf("%s = %s", Name(x), Name(x)))
}

func (s Equal) Symm(e Sort) Sort {
	// symmetric
	panic("not implemented")
}

func (s Equal) Trans(e1 Sort, e2 Sort) Sort {
	// transitive
	panic("not implemented")
}

type equalTerm struct {
	A      Sort
	B      Sort
	Parent Equal
}

func (s equalTerm) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s = %s", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Equal:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}
