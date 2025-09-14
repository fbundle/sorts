package el_sorts

func Eval(sort Sort) Sort {
	// TODO

	switch sort := sort.(type) {
	case Beta:
		return sort
		
	default:
		return sort
	}
}
