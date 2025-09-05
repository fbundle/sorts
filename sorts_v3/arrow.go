package sorts

import "fmt"

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) view() view {
	level := max(Level(s.A), Level(s.B))
	return view{
		view: level,
		name: fmt.Sprintf("%s -> %s", Name(s.A), Name(s.B)),
		parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				return false
			case Arrow:
				// reverse cast for arg
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !LessEqual(d.A, s.A) {
					return false
				}
				// normal cast for body
				return LessEqual(s.B, d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}

func (s Arrow) ModusPonens(a Sort) Sort {
	// take (a: A) give (b: B)
	mustTermOf(a, s.A)
	return NewAtom(
		Level(s.B)-1,
		fmt.Sprintf("(modus_ponens %s %s)", Name(s), Name(a)),
		s.B,
	)
}

func (s Arrow) Intro(f func(Sort) Sort) Sort {
	// take a func that maps (a: A) into (b: B)
	// give (x: A -> B)

	// verify
	a := NewAtom(Level(s.A)-1, "a_dummy", s.A)
	b := f(a)
	mustTermOf(b, s.B)

	// verify ok
	return NewAtom(
		Level(s)-1,
		fmt.Sprintf("(implies_intro %s %p)", Name(s), f),
		s,
	)
}
