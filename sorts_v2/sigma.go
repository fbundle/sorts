package sorts

import (
	"fmt"

	"github.com/fbundle/sorts/adt"
)

func newSigma(ss SortSystem, a InhabitedSort, b DependentSort) adt.Option[Sort] {
	return adt.Some[Sort](sigma{
		a:  a,
		b:  b,
		ss: ss,
	})
}

// sigma - Sigma-type (dependent pair)
// (x: A, y: B(x))
type sigma struct {
	a  InhabitedSort
	b  DependentSort
	ss SortSystem
}

func (s sigma) Level() int {
	return max(s.a.Sort().Level(), s.b.Sort().Level())
}

func (s sigma) Name() string {
	return fmt.Sprintf("(x: %s, y: %s(x))", s.a.Sort().Name(), s.b.Sort().Name())
}

func (s sigma) Parent() InhabitedSort {
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s sigma) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
