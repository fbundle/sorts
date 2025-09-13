package sorts

import "fmt"

type Prod struct {
	H Name
	A Sort
	B Sort
}

func (s Prod) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Prod{A: a.Parent(s.A), B: a.Parent(s.B)}, // smallest type containing A × B
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return a.LessEqual(s.A, d.A) && a.LessEqual(s.B, d.B)
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

	a := NewTerm(s.A, fmt.Sprintf("(left %s)", Name(t)))
	b := NewTerm(s.B, fmt.Sprintf("(right %s)", Name(t)))
	return a, b
}
