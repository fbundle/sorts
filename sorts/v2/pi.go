package sorts

// pi - Pi-type (dependent function)
// (x: Arg) -> Body(x)
type pi struct {
	arg  InhabitedSort
	body func(Sort) Sort
	ss   SortSystem
}
