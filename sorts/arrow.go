package sorts

type Arrow struct {
	Head Name
	A    Sort
	B    Sort
}

func (s Arrow) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.Head, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Arrow{Head: s.Head, A: a.Parent(s.A), B: a.Parent(s.B)},
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
