package sorts

import "fmt"

var _ Sort = Arrow{}

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) View() View {
	aView, bView := s.A.View(), s.B.View()
	level := max(aView.Level, bView.Level)
	return View{
		Level: level,
		Name:  fmt.Sprintf("%s -> %s", aView.Name, bView.Name),
		Parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		LessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				return false
			case Arrow:
				// reverse cast for arg
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !d.A.View().LessEqual(s.A) {
					return false
				}
				// normal cast for body
				return s.B.View().LessEqual(d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}
