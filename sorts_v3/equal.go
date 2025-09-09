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
		parent: s.Parent,
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

func (s Equal) Symm(t equalTerm) equalTerm {
	// symmetric
	MustTermOf(t, s)
	return equalTerm{
		A:      t.B,
		B:      t.A,
		Parent: s,
	}
}

func (s Equal) Trans(t1 equalTerm, t2 equalTerm) Sort {
	// transitive
	MustTermOf(t1, s)
	MustTermOf(t2, s)
	
}
