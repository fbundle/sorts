package sorts

import (
	"fmt"
	"strings"
)

func mustName(err error, form Form) Name {
	if name, ok := form.(Name); ok {
		return name
	}
	panic(err)
}

func parseErr(cmd Name, args []string, suffices ...string) error {
	argStr := strings.Join(args, " ")
	suffixStr := strings.Join(suffices, " ")
	return fmt.Errorf("%s must be (%s %s) %s", cmd, cmd, argStr, suffixStr)
}
