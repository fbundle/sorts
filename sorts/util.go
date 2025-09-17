package sorts

func serialize(s Sort) []Sort {
	if s, ok := s.(Arrow); ok {
		body := serialize(s.B)
		return append([]Sort{s.A}, body...)
	} else {
		return []Sort{s}
	}
}
