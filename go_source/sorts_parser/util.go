package sorts_parser

import (
	"fmt"
	"strings"
)

func makeForm(cmd Name, args ...string) string {
	return fmt.Sprintf("(%s %s)", cmd, strings.Join(args, " "))
}

func compileErr(actual Form, cmd []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	cmdStr := makeForm(Name(cmd[0]), cmd[1:]...)
	return fmt.Errorf("%s must be %s %s got %s", cmd, cmdStr, suffixStr, actual)
}
func mustType[T any](err error, o any) T {
	if v, ok := o.(T); ok {
		return v
	}
	panic(err)
}
