package sorts

// Pi - dependent function type Π_{x: A} B(x)
type Pi struct {
	H Name
	X Sort       // some (x:A)
	B Dept[Sort] // must have the same level for every x: A
}

func (s Pi) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:  List{s.H, a.Form(s.X), a.Form(s.B)},
		level: max(a.Level(a.Parent(s.X)), a.Level(s.B.Apply(s.X))),
	}
}
