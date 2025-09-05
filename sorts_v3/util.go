package sorts

func mustType(x Sort, X Sort) {
	if x.View().Parent.Sort != X {
		panic("type_error")
	}
}
