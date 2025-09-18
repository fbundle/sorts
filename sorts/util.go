package sorts

import (
	"fmt"
	"strings"
)

func mustType[T any](err error, o any) T {
	if v, ok := o.(T); ok {
		return v
	}
	panic(err)
}

func makeForm(cmd Name, args ...string) string {
	return fmt.Sprintf("(%s %s)", cmd, strings.Join(args, " "))
}

func compileErr(cmd Name, args []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be %s %s", cmd, makeForm(cmd, args...), suffixStr)
}

func serialize(s Sort) []Sort {
	if s, ok := s.(Arrow); ok {
		return append([]Sort{s.A}, serialize(s.B)...)
	}
	return []Sort{s}
}

func slicesMap[T1 any, T2 any](input []T1, f func(T1) T2) []T2 {
	output := make([]T2, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
