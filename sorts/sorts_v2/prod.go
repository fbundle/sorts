package sorts

import (
	"fmt"

	"github.com/fbundle/lab_public/lab/go_util/pkg/adt"
)

func newProd(ss SortSystem, a Sort, b Sort) adt.Option[Sort] {
	return adt.Some[Sort](prod{
		a:  a,
		b:  b,
		ss: ss,
	})
}

type prod struct {
	a  Sort
	b  Sort
	ss SortSystem
}

func (s prod) Level() int {
	return max(s.a.Level(), s.b.Level())
}

func (s prod) Name() string {
	return fmt.Sprintf("%s × %s", s.a.Name(), s.b.Name())
}

func (s prod) Parent() InhabitedSort {
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s prod) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
