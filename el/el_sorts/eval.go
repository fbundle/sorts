package el_sorts

func Eval(sort Sort) Sort {

	switch sort := sort.(type) {
	case Beta:
		return sort // TODO
	case Let:
		return sort // TODO
	case Match:
		return sort // TODO

	default:
		return sort
	}
}
