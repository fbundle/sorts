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

func compileErr(actual Form, cmd []string, suffices ...string) error {
	suffixStr := strings.Join(suffices, " ")
	cmdStr := makeForm(Name(cmd[0]), cmd[1:]...)
	return fmt.Errorf("%s must be %s %s got %s", cmd, cmdStr, suffixStr, actual)
}

func compileAnnot(ctx Context, list List) Annot {
	err := compileErr(list, []string{string(AnnotCmd), "name", "type"})
	if len(list) != 2 {
		panic(err)
	}
	return Annot{
		Name: mustType[Name](err, list[0]),
		Type: ctx.Parse(list[1]),
	}
}

const (
	AnnotCmd Name = ":"
)

type Annot struct {
	Name Name
	Type Sort
}

func (a Annot) Form() Form {
	return List{AnnotCmd, a.Name, a.Type.Form()}
}
