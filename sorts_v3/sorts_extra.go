package sorts

// Inhabited - represents a Sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort // underlying sort
	Child Sort // child/term of Sort
}

func (s Inhabited) sortAttr() sortAttr {
	return s.Sort.sortAttr()
}

// Dependent - represent a type B(x) depends on Sort x
type Dependent struct {
	Name  string
	Apply func(Sort) Sort // take x, return B(x)
}

func (d Dependent) nameAttr() string {
	return d.Name
}
