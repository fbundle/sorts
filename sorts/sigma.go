package sorts

import "fmt"

// Sigma - (x: A, y: B(x)) , similar to Prod
// this is the existential quantifier
type Sigma struct {
	A Sort
	B Dependent
}

func (s Sigma) sortAttr() sortAttr {
	x := dummyTerm(s.A, "x")
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		name:   fmt.Sprintf("Σ(x:%s)%s(x)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sigma:
				y := dummyTerm(d.A, "y")
				dBy := d.B.Apply(y)
				return LessEqual(s.A, d.A) && LessEqual(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B(a)) give (t: Σ(x:A)B(x))
func (s Sigma) Intro(a Sort, b Sort) Sort {
	MustTermOf(a, s.A)
	MustTermOf(b, s.B.Apply(a))
	return dummyTerm(s, fmt.Sprintf("(%s, %s)", Name(a), Name(b)))
}

// Elim - take (t: Σ(x:A)B(x)) give (a: A) (b: B(a))
func (s Sigma) Elim(t Sort) (left Sort, right Sort) {
	MustTermOf(t, s)
	a := dummyTerm(s.A, "a")
	b := dummyTerm(s.B.Apply(a), "b")
	return a, b
}
