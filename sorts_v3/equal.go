package sorts

// Equal - A <-> B
type Equal struct {
	A Sort
	B Sort
}

func (s Equal) attr() sortAttr {
	panic("implement me")
}
