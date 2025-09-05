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
	name      string              // every Sort is identified with a Name (string)
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

type Sort interface {
	attr() sortAttr
}

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
	Sort  Sort            // underlying sort
	Apply func(Sort) Sort // take x, return B(x)
}

func (s Dependent) attr() sortAttr {
	return s.Sort.attr()
}
