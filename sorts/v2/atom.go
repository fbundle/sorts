package sorts

type atom struct {
	level  int
	name   string
	parent Sort
}

func (s atom) Level(ss SortSystem) int {
	return s.level
}

func (s atom) Name(ss SortSystem) string {
	return s.name
}

func (s atom) Parent(ss SortSystem) Sort {
	if s.parent != nil {
		return s.parent
	}
	// default parent
	return ss.Default(s.Level(ss) + 1)
}

func (s atom) LessEqual(ss SortSystem, dst Sort) bool {
	switch d := dst.(type) {
	case atom:
		if s.level != d.level {
			return false
		}
		return ss.LessEqual(s.name, d.name)
	case arrow:
		// cannot compare atom and arrow
		return false
	default:
		panic("unreachable")
	}
}
