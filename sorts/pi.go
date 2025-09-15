package sorts

// Pi - dependent function type Π_{x: A} B(x)
type Pi struct {
	H Name
	A Sort
	B Dept[Sort] // must have the same level for every x: A
}

func (s Pi) sortAttr(a SortAttr) sortAttr {
	x := s.A // TODO - some term of A

	return sortAttr{
		form:  List{s.H, a.Form(x), a.Form(s.B)},
		level: max(a.Level(a.Parent(x)), a.Level(s.B.Apply(x))),
	}
}
