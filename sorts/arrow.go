package sorts

import "fmt"

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s -> %s", Name(s.A), Name(s.B)),
		level:  level,
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
	termName := fmt.Sprintf("(%s %s)", Name(f), Name(a))
	return makeTerm(termName, s.B)
}

// Intro - take a go function (repr) that maps (a: A) into (b: B)  give (f: A -> B)
func (s Arrow) Intro(name string, repr func(Sort) Sort) Sort {
	// verify
	a := makeTerm("a", s.A) // dummy term
	b := repr(a)
	mustTermOf(b, s.B)

	// verify ok
	return makeTerm(name, s)
}
