package sorts

// Sigma - (x: A, y: B(x)) , similar to Prod
// this is the existential quantifier
type Sigma struct {
	A Inhabited
	B Dependent
}

func (s Sigma) attr() sortAttr {
	x := dummyTerm(s.A, "x")
	level := max(Level(s.A), Level(s.B(x)))
	return sortAttr{
		level: level,
	}
}
