package sorts

import (
	"fmt"
)

func ParseListArrow(H Name) ParseListFunc {
	return func(parse ParseFunc, list List) (Sort, error) {
		if len(list) != 3 {
			return nil, fmt.Errorf("arrow must be %s domain codomain", H)
		}
		if list[0] != H {
			return nil, fmt.Errorf("arrow must be %s domain codomain", H)
		}
		a, err := parse(list[1])
		if err != nil {
			return nil, err
		}
		b, err := parse(list[2])
		if err != nil {
			return nil, err
		}
		return Arrow{H: H, A: a, B: b}, nil
	}
}

type Arrow struct {
	H Name
	A Sort
	B Sort
}

func (s Arrow) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Arrow{H: s.H, A: a.Parent(s.A), B: a.Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !a.LessEqual(d.A, s.A) {
					return false
				}
				return a.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}
