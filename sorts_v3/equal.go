package sorts

import "fmt"

// Equal - A <-> B
type Equal struct {
	A Sort
	B Sort
}

func (s Equal) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s <-> %s", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Equal:
				return LessEqual(Arrow{s.A, s.B}, Arrow{d.A, d.B}) && LessEqual(Arrow{d.A, d.B}, Arrow{s.A, s.B})
			default:
				return false
			}
		},
	}
}

func (s Equal) Elim(t Sort) (Sort, Sort) {
	MustTermOf(t, s)
	l2r := Arrow{s.A, s.B}
	r2l := Arrow{s.B, s.A}
	return dummyTerm(l2r, "l2r"), dummyTerm(r2l, "rl2")
}

func (s Equal) Intro(name string, l2rFunc func(Sort) Sort, r2lFunc func(Sort) Sort) Sort {
	l2r := Arrow{s.A, s.B}
	r2l := Arrow{s.B, s.A}
	l2r.Intro(name, l2rFunc)
	r2l.Intro(name, r2lFunc)

	return dummyTerm(s, name)
}
