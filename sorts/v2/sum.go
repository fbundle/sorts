package sorts

import (
	"fmt"

	"github.com/fbundle/sorts/adt"
)

func newSum(ss SortSystem, a Sort, b Sort) adt.Option[Sort] {
	return adt.Some[Sort](sum{
		a:  a,
		b:  b,
		ss: ss,
	})
}

type sum struct {
	a  Sort
	b  Sort
	ss SortSystem
}

func (s sum) Level() int {
	return max(s.a.Level(), s.b.Level())
}

func (s sum) Name() string {
	return fmt.Sprintf("%s + %s", s.a.Name(), s.b.Name())
}

func (s sum) Parent() InhabitedSort {
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s sum) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
