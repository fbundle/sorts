package sorts

// Inhabited - represents a Sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort // underlying sort
	Child Sort // child/term of Sort
}

func (s Inhabited) attr() sortAttr {
	return s.Sort.attr()
}

// Dependent - represent a type B(x) depends on Sort x
type Dependent struct {
	Name  string
	Apply func(Sort) Sort // take x, return B(x)
}

func (d Dependent) attr() sortAttr {
	// Dependent is not a sort so it doesn't have all attributes
	return sortAttr{
		name: d.Name,
	}
}
