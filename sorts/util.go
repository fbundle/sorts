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

// compileErr - TODO change function signature into compileErr(actual Form, suffices ...string)
// TODO - suffices will contains command (from makeForm) and other suffices
func compileErr(actual Form, cmd Name, args []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be %s %s got %s", cmd, makeForm(cmd, args...), suffixStr, actual)
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

func slicesReduce[T1 any, T2 any](input []T1, init T2, f func(T2, T1) T2) T2 {
	output := init
	for _, v := range input {
		output = f(output, v)
	}
	return output
}

func slicesReverse[T any](input []T) []T {
	output := make([]T, len(input))
	for i, v := range input {
		output[len(input)-1-i] = v
	}
	return output
}

func slicesForEach[T any](input []T, f func(T)) {
	for _, v := range input {
		f(v)
	}
}
