package sorts

import (
	"github.com/fbundle/sorts/form"
)

const (
	ProdName      form.Name = "⊗"
	ProdLeftName  form.Name = "left"
	ProdRightName form.Name = "right"
)

func mustParseProd(parse mustParseFunc, args form.List) Sort {
	if len(args) != 2 {
		panic(typeErr)
	}
	A := parse(args[0])
	B := parse(args[1])
	return Prod{A, B}
}

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		repr:   form.List{ProdName, Repr(s.A), Repr(s.B)},
		level:  level,
		parent: Prod{A: Parent(s.A), B: Parent(s.B)}, // smallest type containing A × B
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return SubTypeOf(s.A, d.A) && SubTypeOf(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B) give (a, b): A × B
func (s Prod) Intro(a Sort, b Sort) Sort {
	mustTermOf(a, s.A)
	mustTermOf(b, s.B)
	return Prod{A: a, B: b}
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(t Sort) (left Sort, right Sort) {
	mustTermOf(t, s)
	if t, ok := t.(Prod); ok {
		return t.A, t.B
	}

	a := newAtomTerm(form.List{ProdLeftName, Repr(s)}, s.A)
	b := newAtomTerm(form.List{ProdRightName, Repr(s)}, s.B)
	return a, b
}
