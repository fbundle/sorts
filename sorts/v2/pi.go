package sorts

// pi - Pi-type (dependent function)
// (x: Arg) -> Body(x)
type pi struct {
	arg  Sort
	body func(Sort) Sort
}
