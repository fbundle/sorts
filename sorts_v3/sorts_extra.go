package sorts

// Inhabited - represents a WithSort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  WithSort // underlying sort
	Child WithSort // child/term of Sort
}

func (s Inhabited) attr() sortAttr {
	return s.Sort.attr()
}

// Dependent - represent a type B(x) depends on WithSort x
type Dependent struct {
	Name  string
	Apply func(WithSort) WithSort // take x, return B(x)
}

func (d Dependent) attr() nameAttr {
	// Dependent is not a sort so it doesn't have all attributes
	return nameAttr{
		name: d.Name,
	}
}
