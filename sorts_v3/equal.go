package sorts

// Equal - A <-> B
type Equal struct {
	A WithSort
	B WithSort
}

func (s Equal) attr() sortAttr {
	panic("implement me")
}
