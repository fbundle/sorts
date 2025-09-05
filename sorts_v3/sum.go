package sorts

import "fmt"

var _ Sort = Sum{}

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) view() view {
	level := max(Level(s.A), Level(s.B))
	return view{
		view: level,
		name: fmt.Sprintf("%s + %s", Name(s.A), Name(s.B)),
		parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}

func (s Sum) IntroLeft(a Sort) Sort {
	// take (a: A) give (x: A + B)
	mustTermOf(a, s.A)
	return NewAtom(
		Level(s)-1,
		fmt.Sprintf("(intro_left %s %s)", Name(s), Name(a)),
		s,
	)
}

func (s Sum) IntroRight(b Sort) Sort {
	// take (b: B) give (x: A + B)
	mustTermOf(b, s.B)
	return NewAtom(
		Level(s)-1,
		fmt.Sprintf("(intro_right %s %s)", Name(s), Name(b)),
		s,
	)
}

func (s Sum) ByCases(h1 Arrow, h2 Arrow) Sort {
	// take (h1: A -> X) (h2: B -> X) give (x: X)
	mustSubType(h1.A, s.A)
	mustSubType(h2.A, s.B)
	mustSubType(h1.B, h2.B)
	X := h1.B
	return NewAtom(
		Level(X)-1,
		fmt.Sprintf("(by_cases %s %s %s)", Name(s), Name(h1), Name(h2)),
		X,
	)
}
