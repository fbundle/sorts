package sorts

import "fmt"

func serialize(s Sort) []Sort {
	if s, ok := s.(Arrow); ok {
		body := serialize(s.B)
		return append([]Sort{s.A}, body...)
	} else {
		return []Sort{s}
	}
}

func deserialize(s []Sort) Sort {
	if len(s) == 0 {
		panic(fmt.Errorf("empty sort array"))
	}
	output := s[len(s)-1]
	for i := len(s) - 2; i >= 0; i-- {
		output = Arrow{s[i], output}
	}
	return output
}
