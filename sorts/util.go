package sorts

import (
	"fmt"
	"strings"
)

// serialize - turn type Arrow(a, Arrow(b, c)) into []Sort{a, b, c}
func serialize(s Sort) []Sort {
	if s, ok := s.(Arrow); ok {
		body := serialize(s.B)
		return append([]Sort{s.A}, body...)
	} else {
		return []Sort{s}
	}
}

// deserialize - turn []Sort{a, b, c} into Arrow(a, Arrow(b, c))
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

func parseErr(cmd Name, args []string, suffices ...string) error {
	argStr := strings.Join(args, " ")
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be (%s %s) %s", cmd, cmd, argStr, suffixStr)
}
