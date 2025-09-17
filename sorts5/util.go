package sorts5

func must[T1 any, T2 any](err error, pred func(T1) (T2, bool)) func(t1 T1) T2 {
	return func(t1 T1) T2 {
		t2, ok := pred(t1)
		if !ok {
			panic(err)
		}
		return t2
	}
}

func isName(form Form) (Name, bool) {
	name, ok := form.(Name)
	return name, ok
}
func isList(form Form) (List, bool) {
	list, ok := form.(List)
	return list, ok
}
