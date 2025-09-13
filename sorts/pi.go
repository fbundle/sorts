package sorts

import (
	"github.com/fbundle/sorts/form"
)

const (
	PiName form.Name = "Π"
)

// Pi - (x: A) -> (y: B(x)) similar to Arrow
// this is the universal quantifier
type Pi struct {
	A Sort
	B Dependent
}

func (s Pi) sortAttr() sortAttr {
	x := newTerm(form.Name("x"), s.A) // dummy term
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		repr:   form.List{PiName, Repr(s.A), Repr(s.B)},
		level:  level,
		parent: nil, // TODO
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Pi:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				if !SubTypeOf(d.A, s.A) {
					return false
				}
				y := newTerm(form.Name("y"), d.A)
				dBy := d.B.Apply(y)
				return SubTypeOf(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take a func that maps (a: A) into (b: B(a))  give (f: Π(x:A)B(x))
func (s Pi) Intro(repr form.Form, f func(Sort) Sort) Sort {
	// verify
	a := newTerm(form.Name("a"), s.A)
	b := f(a)
	sBa := s.B.Apply(a)
	mustTermOf(b, sBa) // TODO - think, shouldn't it have to check for every a of type A?

	return newTerm(repr, s)
}

// Elim - take (f: Π(x:A)B(x)) (a: A) give (b: B(a)) - Modus Ponens
func (s Pi) Elim(f Sort, a Sort) Sort {
	mustTermOf(f, s)
	mustTermOf(a, s.A)
	Ba := s.B.Apply(a)

	return newTerm(form.List{Repr(f), Repr(a)}, Ba)
}
