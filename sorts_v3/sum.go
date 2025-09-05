package sorts

import "fmt"

var _ Sort = Sum{}

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) View() View {
	aView, bView := s.A.View(), s.B.View()
	level := max(aView.Level, bView.Level)
	return View{
		Level: level,
		Name:  fmt.Sprintf("%s + %s", aView.Name, bView.Name),
		Parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		LessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return aView.LessEqual(d.A) && bView.LessEqual(d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}

func (s Sum) IntroLeft(a Sort) Sort {
	// take (a: A) give (x: A + B)

}
