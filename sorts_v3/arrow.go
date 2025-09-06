package sorts

import "fmt"

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) attr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("%s -> %s", Name(s.A), Name(s.B)),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !LessEqual(d.A, s.A) {
					return false
				}
				return LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Elim - take (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(a Sort) Sort {
	mustTermOf(a, s.A)
	return dummyTerm(s.B, fmt.Sprintf("(modus_ponens %s %s)", Name(s), Name(a)))
}

// Intro - take a func that maps (a: A) into (b: B)  give (x: A -> B)
func (s Arrow) Intro(f func(Sort) Sort) Sort {
	// verify
	a := dummyTerm(s.A, "a")
	b := f(a)
	mustTermOf(b, s.B)

	// verify ok
	return dummyTerm(s, fmt.Sprintf("(implies_intro %s %p)", Name(s), f))
}
