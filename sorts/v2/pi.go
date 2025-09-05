package sorts

import "github.com/fbundle/sorts/adt"

func newPi(ss SortSystem, arg InhabitedSort, body func(Sort) Sort) adt.Option[Sort] {
	return adt.Some[Sort](pi{
		arg:  arg,
		body: body,
		ss:   ss,
	})
}

// pi - Pi-type (dependent function)
// (x: Arg) -> Body(x)
type pi struct {
	arg  InhabitedSort
	body func(Sort) Sort
	ss   SortSystem
}

func (s pi) Level() int {
	child := s.arg.Child()
	return max(s.arg.Level(), s.body(child).Level())
}

func (s pi) Name() string {
	//TODO implement me
	panic("implement me")
}

func (s pi) Parent() InhabitedSort {
	//TODO implement me
	panic("implement me")
}

func (s pi) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
