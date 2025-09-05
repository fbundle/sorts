package sorts

import "fmt"

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) view() view {
	level := max(Level(s.A), Level(s.B))
	return view{
		view: level,
		name: fmt.Sprintf("%s × %s", Name(s.A), Name(s.B)),
		parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}

// Intro - take (a: A) (b: B) give (a, b): A × B
func (s Prod) Intro(a Sort, b Sort) Sort {
	mustTermOf(a, s.A)
	mustTermOf(b, s.B)
	return NewAtom(
		Level(s)-1,
		fmt.Sprintf("(%s, %s)", Name(a), Name(b)),
		s,
	)
}

// Left - take (x: A × B) give (a: A)
func (s Prod) Left(x Sort) Sort {
	mustTermOf(x, s)
	return NewAtom(
		Level(s.A)-1,
		fmt.Sprintf("(left %s %s)", Name(s), Name(x)),
		s.A,
	)
}

// Right - take (x: A × B) give (b: B)
func (s Prod) Right(x Sort) Sort {
	mustTermOf(x, s)
	return NewAtom(
		Level(s.B)-1,
		fmt.Sprintf("(right %s %s)", Name(s), Name(x)),
		s.B,
	)
}
