package sorts

/*
const (
	EqualTerm form.Name = "Eq"
)

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
	return sortAttr{
		name:   form.List{EqualTerm, Repr(s.A), Repr(s.B)},
		level:  max(Level(s.A), Level(s.B)),
		parent: s.Parent,
		lessEqual: func(dst Sort) bool {
			return false // I haven't find any useful thing for this
		},
	}
}

func (s Equal) sortAttr() sortAttr {
	return sortAttr{
		name:   form.List{EqualTerm, Repr(s.A), Repr(s.B)},
		level:  max(Level(s.A), Level(s.B)),
		parent: nil,
		lessEqual: func(dst Sort) bool {
			return false // I haven't find any useful thing for this
		},
	}
}

func (s Equal) Refl(x Sort) Sort {
	// reflexive
	mustTermOf(x, s.A)
	mustTermOf(x, s.B)
	return equalTerm{A: x, B: x, Parent: s}
}

func (s Equal) Symm(t equalTerm) equalTerm {
	// symmetric
	mustTermOf(t, s)
	return equalTerm{
		A:      t.B,
		B:      t.A,
		Parent: s,
	}
}

func (s Equal) Trans(t1 equalTerm, t2 equalTerm) Sort {
	// transitive
	mustTermOf(t1, s)
	mustTermOf(t2, s)
	if t1.B != t2.A {
		panic("type_error")
	}
	return equalTerm{
		A:      t1.A,
		B:      t2.B,
		Parent: s,
	}
}

*/
