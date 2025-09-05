package sorts

func Name(s Sort) string {
	return s.view().name
}

func Level(s Sort) int {
	return s.view().view
}

func Parent(s Sort) Sort {
	return s.view().parent.Sort
}

func LessEqual(x Sort, y Sort) bool {
	return x.view().lessEqual(y)
}

type view struct {
	view      int                 // universe Level
	name      string              // every sort is identified with a Name (string)
	parent    Inhabited           // (or Type) every sort must have a Parent
	lessEqual func(dst Sort) bool // partial order on sorts
}

type Sort interface {
	view() view
}

// Inhabited - represents a sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort
	Child Sort // (or Term)
}

func (s Inhabited) view() view {
	return s.Sort.view()
}

// Dependent - represent a type B(x) depends on sort x
type Dependent struct {
	Sort  Sort
	Apply func(Sort) Sort // take x, return B(x)
}

func (s Dependent) view() view {
	return s.Sort.view()
}
