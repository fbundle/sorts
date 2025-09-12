package sorts

import "fmt"

// Sigma - (x: A, y: B(x)) , similar to Prod
// this is the existential quantifier
type Sigma struct {
	A Sort
	B Dependent
}

func (s Sigma) sortAttr() sortAttr {
	x := NewTerm(s.A, "x")
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		name:   fmt.Sprintf("Σ(x:%s)%s(x)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sigma:
				y := NewTerm(d.A, "y")
				dBy := d.B.Apply(y)
				return SubTypeOf(s.A, d.A) && SubTypeOf(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B(a)) give (t: Σ(x:A)B(x))
func (s Sigma) Intro(a Sort, b Sort) Sort {
	mustTermOf(a, s.A)
	mustTermOf(b, s.B.Apply(a))
	return NewTerm(s, fmt.Sprintf("(%s, %s)", Name(a), Name(b)))
}

// Elim - take (t: Σ(x:A)B(x)) give (a: A) (b: B(a))
func (s Sigma) Elim(t Sort) (left Sort, right Sort) {
	mustTermOf(t, s)
	a := NewTerm(s.A, "a")
	b := NewTerm(s.B.Apply(a), "b")
	return a, b
}
