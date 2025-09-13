package sorts

import (
	"github.com/fbundle/sorts/form"
)

const (
	SigmaName form.Name = "Σ"
)

// Sigma - (x: A, y: B(x)) , similar to Prod
// this is the existential quantifier
type Sigma struct {
	A Sort
	B Dependent
}

func (s Sigma) sortAttr() sortAttr {
	x := newAtomTerm(form.Name("x"), s.A)
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		repr:   form.List{SigmaName, Repr(s.A), Repr(s.B)},
		level:  level,
		parent: nil,
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sigma:
				y := newAtomTerm(form.Name("y"), d.A)
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
	return newAtomTerm(form.List{SigmaName, Repr(a), Repr(b)}, s)
}

// Elim - take (t: Σ(x:A)B(x)) give (a: A) (b: B(a))
func (s Sigma) Elim(t Sort) (left Sort, right Sort) {
	mustTermOf(t, s)

	a := newAtomTerm(form.Name("a"), s.A)
	b := newAtomTerm(form.Name("b"), s.B.Apply(a))
	return a, b
}
