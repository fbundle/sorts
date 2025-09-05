package sorts

import "fmt"

var _ Sort = Prod{}

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) View() View {
	aView, bView := s.A.View(), s.B.View()
	level := max(aView.Level, bView.Level)
	return View{
		Level: level,
		Name:  fmt.Sprintf("%s × %s", aView.Name, bView.Name),
		Parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		LessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return aView.LessEqual(d.A) && bView.LessEqual(d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}

func (s Prod) Intro(a Sort, b Sort) Sort {
	// take (a: A) (b: B) give (a, b): A × B
	mustType(a, s.A)
	mustType(b, s.B)
	return NewAtom(
		s.View().Level-1,
		fmt.Sprintf("(%s, %s)", a.View().Name, b.View().Name),
		s,
	)
}

func (s Prod) Left(x Sort) Sort {
	// take (x: A × B) give (a: A)
	mustType(x, s)
	return NewAtom(
		s.A.View().Level-1,
		fmt.Sprintf("(left %s)", x.View().Name),
		s.A,
	)
}
func (s Prod) Right(x Sort) Sort {
	// take (x: A × B) give (b: B)
	mustType(x, s)
	return NewAtom(
		s.B.View().Level-1,
		fmt.Sprintf("(right %s)", x.View().Name),
		s.B,
	)
}
