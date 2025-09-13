package sorts

import (
	"fmt"
	"strings"
)

func mustTermOf(x Sort, X Sort) {
	mustTypeOf(TermOf(x, X))
}

func mustTypeOf(ok bool) {
	if !ok {
		panic("type_error")
	}
}

func (n Node[T]) String() string {
	nameList := []string{fmt.Sprint(n.Value)}
	for _, child := range n.Children {
		nameList = append(nameList, child.String())
	}
	return "(" + strings.Join(nameList, " ") + ")"
}
