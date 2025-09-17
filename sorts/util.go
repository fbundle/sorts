package sorts

import (
	"fmt"
	"strings"
)

// serialize - turn type Arrow(a, Arrow(b, c)) into []Sort1{a, b, c}
func serialize(s Sort1) []Sort1 {
	if s, ok := s.(Arrow); ok {
		body := serialize(s.B)
		return append([]Sort1{s.A}, body...)
	} else {
		return []Sort1{s}
	}
}

// deserialize - turn []Sort1{a, b, c} into Arrow(a, Arrow(b, c))
func deserialize(s []Sort1) Sort1 {
	if len(s) == 0 {
		panic(fmt.Errorf("empty sort array"))
	}
	output := s[len(s)-1]
	for i := len(s) - 2; i >= 0; i-- {
		output = Arrow{s[i], output}
	}
	return output
}

func parseErr(cmd Name, args []string, suffices ...string) error {
	argStr := strings.Join(args, " ")
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be (%s %s) %s", cmd, cmd, argStr, suffixStr)
}
