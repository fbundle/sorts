package sorts

type SortView struct {
	Level     int                 // universe level
	Name      string              // every sort is identified with a name (string)
	Parent    Inhabited           // (or Type) every sort must have a parent
	LessEqual func(dst Sort) bool // partial order on sorts
}

type Sort interface {
	View() SortView
}

// Inhabited - represents a sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort
	Child Sort // (or Term)
}

func (s Inhabited) View() SortView {
	return s.Sort.View()
}

// Dependent - represent a type B(x) depends on sort x
type Dependent struct {
	Sort  Sort
	Apply func(Sort) Sort // take x, return B(x)
}

func (s Dependent) View() SortView {
	return s.Sort.View()
}
