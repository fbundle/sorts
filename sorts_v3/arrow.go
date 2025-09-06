package sorts

import "fmt"

// Arrow - (A -> B)
type Arrow struct {
	A WithSort
	B WithSort
}

func (s Arrow) nameAttr() string {
	return fmt.Sprintf("%s -> %s", Name(s.A), Name(s.B))
}

func (s Arrow) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst WithSort) bool {
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

// Elim - take (f: A -> B) (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(f WithSort, a WithSort) WithSort {
	mustTermOf(f, s)
	mustTermOf(a, s.A)
	return dummyTerm(s.B, fmt.Sprintf("(elim %s %s)", Name(s), Name(a)))
}

// Intro - take a func that maps (a: A) into (b: B)  give (x: A -> B)
func (s Arrow) Intro(f func(WithSort) WithSort) WithSort {
	// verify
	a := dummyTerm(s.A, "a")
	b := f(a)
	mustTermOf(b, s.B)

	// verify ok
	return dummyTerm(s, fmt.Sprintf("(implies_intro %s %p)", Name(s), f))
}
