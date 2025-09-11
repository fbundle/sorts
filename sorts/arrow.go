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

// Elim - take (f: A -> B) (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(arrow Sort, a Sort) Sort {
	MustTermOf(arrow, s)
	MustTermOf(a, s.A)
	return NewTerm(s.B, fmt.Sprintf("(%s %s)", Name(arrow), Name(a)))
}

// Intro - take a func that maps (a: A) into (b: B)  give (x: A -> B)
func (s Arrow) Intro(name string, arrow func(Sort) Sort) Sort {
	// verify
	a := NewTerm(s.A, "a")
	b := arrow(a)
	MustTermOf(b, s.B)

	// verify ok
	return NewTerm(s, name)
}
