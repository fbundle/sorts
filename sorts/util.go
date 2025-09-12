package sorts

func mustTermOf(x Sort, X Sort) {
	mustTypeOf(TermOf(x, X))
}

func mustTypeOf(ok bool) {
	if !ok {
		panic("type_error")
	}
}
