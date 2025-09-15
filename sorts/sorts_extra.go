package sorts

// Inhabited - represents a Sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort // underlying sort
	Child Sort
}

func (s Inhabited) sortAttr(a SortAttr) sortAttr {
	return s.Sort.sortAttr(a)
}

// Dept - represent a type B(x) depends on Sort x
type Dept[T any] struct {
	Form  Form
	Apply func(T) Sort // take x, return B(x)
}

// Inductive - inductive type
type Inductive interface {
	Sort
}
