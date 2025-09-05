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
	aView, bView := a.View(), b.View()
	if aView.Parent.Sort != s.A || bView.Parent.Sort != s.B {
		panic("type_error")
	}
	// take (a: A) (b: B) give (a, b): A × B
	return NewAtom(
		s.View().Level-1,
		fmt.Sprintf("(%s, %s)", aView.Name, bView.Name),
		s,
	)
}
