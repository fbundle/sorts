package sorts

import (
	"fmt"
	"strings"
	"sync/atomic"
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

func compileErr(actual Form, cmd []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	cmdStr := makeForm(Name(cmd[0]), cmd[1:]...)
	return fmt.Errorf("%s must be %s %s got %s", cmd, cmdStr, suffixStr, actual)
}

var count uint64 = 0

func nextCount() uint64 {
	return atomic.AddUint64(&count, 1)
}

type code struct{}

func (c code) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (c code) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (c code) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}
