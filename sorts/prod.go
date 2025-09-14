package sorts

import "fmt"

func ListCompileProd(H Name) ListCompileFunc {
	return func(parse func(form Form) Sort, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("prod must be %s A B", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("prod must be %s A B", H))
		}
		return Prod{H: H, A: parse(list[1]), B: parse(list[2])}
	}
}

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
func (s Prod) Intro(sa SortAttr, a Sort, b Sort) Sort {
	must(sa).termOf(a, s.A)
	must(sa).termOf(b, s.B)
	return Prod{A: a, B: b}
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(sa SortAttr, t Sort) (left Sort, right Sort) {
	must(sa).termOf(t, s)
	if t, ok := t.(Prod); ok {
		return t.A, t.B
	}

	a := NewAtomTerm(sa, List{Name("left"), sa.Form(t)}, s.A)
	b := NewAtomTerm(sa, List{Name("right"), sa.Form(t)}, s.B)
	return a, b
}
