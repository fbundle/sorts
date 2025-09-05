package sorts

import "fmt"

var _ Sort = Arrow{}

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) view() view {
	level := max(Level(s.A), Level(s.B))
	return view{
		view: level,
		name: fmt.Sprintf("%s -> %s", Name(s.A), Name(s.B)),
		parent: Inhabited{
			Sort:  defaultSort(nil, level+1),
			Child: s,
		},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				return false
			case Arrow:
				// reverse cast for arg
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !LessEqual(d.A, s.A) {
					return false
				}
				// normal cast for body
				return LessEqual(s.B, d.B)
			default:
				panic("type_error - should catch all types")

			}
		},
	}
}
