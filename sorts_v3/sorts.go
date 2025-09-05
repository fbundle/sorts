package sorts

func Name(s Sort) string {
	return s.attr().name
}

func Level(s Sort) int {
	return s.attr().view
}

func Parent(s Sort) Sort {
	return s.attr().parent
}

func LessEqual(x Sort, y Sort) bool {
	return x.attr().lessEqual(y)
}

type sortAttr struct {
	view      int                 // universe Level
	name      string              // every sort is identified with a Name (string)
	parent    Sort                // (or Type) every sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

type Sort interface {
	attr() sortAttr
}

// Inhabited - represents a sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort
	Child Sort // (or Term)
}

func (s Inhabited) attr() sortAttr {
	return s.Sort.attr()
}

// Dependent - represent a type B(x) depends on sort x
type Dependent struct {
	Sort  Sort
	Apply func(Sort) Sort // take x, return B(x)
}

func (s Dependent) attr() sortAttr {
	return s.Sort.attr()
}
