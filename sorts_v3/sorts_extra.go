package sorts

// Inhabited - represents a WithSort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  WithSort // underlying sort
	Child WithSort // child/term of Sort
}

func (s Inhabited) nameAttr() string {
	return Name(s.Sort)
}

func (s Inhabited) sortAttr() sortAttr {
	return s.Sort.sortAttr()
}

// Dependent - represent a type B(x) depends on WithSort x
type Dependent struct {
	Name  string
	Apply func(WithSort) WithSort // take x, return B(x)
}

func (d Dependent) nameAttr() string {
	return d.Name
}
