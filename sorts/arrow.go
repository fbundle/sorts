package sorts

import (
	"fmt"
)

func ListParseArrow(H Name) ListParseFunc {
	return func(parse ParseFunc, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("arrow must be %s domain codomain", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("arrow must be %s domain codomain", H))
		}
		a := parse(list[1])
		b := parse(list[2])
		return Arrow{H: H, A: a, B: b}
	}
}

type Arrow struct {
	H Name
	A Sort
	B Sort
}

func (s Arrow) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Arrow{H: s.H, A: a.Parent(s.A), B: a.Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !a.LessEqual(d.A, s.A) {
					return false
				}
				return a.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Elim - take (f: A -> B) (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(sa SortAttr, f Sort, a Sort) Sort {
	must(sa).termOf(f, s)
	must(sa).termOf(a, s.A)

	name := Name(fmt.Sprintf("elim_%s_%s", sa.Form(f), sa.Form(a)))
	return NewAtomTerm(sa, name, s.B)
}

// Intro - take a func that maps (a: A) into (b: B)  give (f: A -> B)
func (s Arrow) Intro(sa SortAttr, name Name, f func(Sort) Sort) Sort {
	// verify
	a := NewAtomTerm(sa, "a", s.A)
	b := f(a)

	must(sa).termOf(b, s.B)

	// verify ok

	return NewAtomTerm(sa, name, s)
}
