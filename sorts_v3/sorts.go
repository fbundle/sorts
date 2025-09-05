package sorts

type SortView struct {
	Level     int                 // universe level
	Name      string              // every sort is identified with a name (string)
	Parent    InhabitedSort       // (or Type) every sort must have a parent
	LessEqual func(dst Sort) bool // partial order on sorts
}

type Sort interface {
	View() SortView
}

// InhabitedSort - represents a sort with at least one child
// (true theorems have proofs)
type InhabitedSort struct {
	Sort  Sort
	Child Sort // (or Term)
}

func (s InhabitedSort) View() SortView {
	return s.Sort.View()
}

// DependentSort - represent a type B(x) depends on sort x
type DependentSort struct {
	Sort  Sort
	Apply func(Sort) Sort // take x, return B(x)
}

func (d DependentSort) View() SortView {
	return d.Sort.View()
}
