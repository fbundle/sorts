package sorts

import (
	"github.com/fbundle/sorts/form"
)

const (
	ArrowName form.Name = "Arrow"
)

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) sortAttr() sortAttr {
	return sortAttr{
		repr:   form.List{ArrowName, Repr(s.A), Repr(s.B)},
		level:  max(Level(s.A), Level(s.B)),
		parent: Arrow{A: Parent(s.A), B: Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !SubTypeOf(d.A, s.A) {
					return false
				}
				return SubTypeOf(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Elim - take (f: A -> B) (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(f Sort, a Sort) Sort {
	mustTermOf(f, s)
	mustTermOf(a, s.A)
	return NewTerm(form.List{Repr(f), Repr(a)}, s.B)
}

// Intro - take a go function (repr) that maps (a: A) into (b: B)  give (f: A -> B)
func (s Arrow) Intro(repr form.Form, f func(Sort) Sort) Sort {
	// verify
	a := NewTerm(form.Name("a"), s.A) // dummy term
	b := f(a)
	mustTermOf(b, s.B)

	// verify ok
	return NewTerm(repr, s)
}
