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

func parseErr(cmd Name, args []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be %s %s", cmd, makeForm(cmd, args...), suffixStr)
}
