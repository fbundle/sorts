package sorts

// Equal - A <-> B
type Equal struct {
	A WithSort
	B WithSort
}

func (s Equal) nameAttr() string {
	panic("implement me")
}

func (s Equal) sortAttr() sortAttr {
	panic("implement me")
}
