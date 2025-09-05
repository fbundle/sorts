package sorts

import (
	"fmt"

	"github.com/fbundle/lab_public/lab/go_util/pkg/adt"
)

func newPi(ss SortSystem, arg InhabitedSort, body DependentSort) adt.Option[Sort] {
	return adt.Some[Sort](pi{
		arg:  arg,
		body: body,
		ss:   ss,
	})
}

// pi - Pi-type (dependent function)
// (x: A) -> B(x)
type pi struct {
	arg  InhabitedSort
	body DependentSort
	ss   SortSystem
}

func (s pi) Level() int {
	return max(s.arg.Sort().Level(), s.body.Sort().Level())
}

func (s pi) Name() string {
	return fmt.Sprintf("(x: %s) -> %s(x)", s.arg.Sort().Name(), s.body.Sort().Name())
}

func (s pi) Parent() InhabitedSort {
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s pi) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
