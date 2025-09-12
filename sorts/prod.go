package sorts

import "fmt"

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s × %s", Name(s.A), Name(s.B)),
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
	MustTermOf(a, s.A)
	MustTermOf(b, s.B)
	return Prod{A: a, B: b}
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(t Sort) (left Sort, right Sort) {
	MustTermOf(t, s)
	if t, ok := t.(Prod); ok {
		return t.A, t.B
	}

	a := NewTerm(s.A, fmt.Sprintf("(left %s)", Name(t)))
	b := NewTerm(s.B, fmt.Sprintf("(right %s)", Name(t)))
	return a, b
}
