package sorts

import "fmt"

type Equal struct {
	A Sort
	B Sort
}

type equalTerm struct {
	A      Sort
	B      Sort
	Parent Equal
}

func (s equalTerm) sortAttr() sortAttr {
	name := fmt.Sprintf("%s = %s", Name(s.A), Name(s.B))
	if nameWithType {
		name = fmt.Sprintf("(%s : %s)", name, Name(s.Parent))
	}
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   name,
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			return false // I haven't find any useful thing for this
		},
	}
}

func (s Equal) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("Eq(%s, %s)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			return false // I haven't find any useful thing for this
		},
	}
}

func (s Equal) Refl(x Sort) Sort {
	// reflexive
	MustTermOf(x, s.A)
	MustTermOf(x, s.B)
	return equalTerm{A: x, B: x, Parent: s}
}

func (s Equal) Symm(e Sort) Sort {
	// symmetric
	panic("not implemented")
}

func (s Equal) Trans(e1 Sort, e2 Sort) Sort {
	// transitive
	panic("not implemented")
}
