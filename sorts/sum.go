package sorts

import "fmt"

func ListCompileSum(H Name) ListCompileFunc {
	return func(compile func(form Form) Sort, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("sum must be %s A B", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("sum must be %s A B", H))
		}
		return Sum{H: H, A: compile(list[1]), B: compile(list[2])}
	}
}

type Sum struct {
	H Name
	A Sort
	B Sort
}

func (s Sum) sortAttr(sa SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, sa.Form(s.A), sa.Form(s.B)},
		level:  max(sa.Level(s.A), sa.Level(s.B)),
		parent: Sum{A: sa.Parent(s.A), B: sa.Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return sa.LessEqual(s.A, d.A) && sa.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - IntroLeft or IntroRight
func (s Sum) Intro(sa SortAttr, a Sort, b Sort) Sort {
	if a != nil {
		// IntroLeft - take (a: A) give (x: A + B)
		must(sa).termOf(a, s.A)
		return a
	} else {
		// IntroRight - take (b: B) give (x: A + B)
		must(sa).termOf(b, s.B)
		return b
	}
}

// ByCases - take (t: A + B) (h1: A -> X) (h2: B -> X) give (x: X)
func (s Sum) ByCases(sa SortAttr, t Sort, h1 Sort, h2 Sort) Sort {
	must(sa).termOf(t, s)
	X := sa.Parent(h1).(Arrow).B
	must(sa).termOf(h1, Arrow{s.H, s.A, X})
	must(sa).termOf(h2, Arrow{s.H, s.B, X})

	return NewAtomTerm(sa, List{Name("by_cases"), sa.Form(t), sa.Form(h1), sa.Form(h2)}, X)
}
