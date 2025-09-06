package sorts

// Pi - (x: A) -> (y: B(x))
type Pi struct {
	A Sort // must be inhabited // TODO consider using Inhabited
	B Dependent
}

func (s Pi) attr() sortAttr {
	return sortAttr{}
}
