package el_typesafe

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type _partialObject struct {
	_sort sorts.Sort
	_next Expr
}

func newPartialObject(next Expr) _partialObject {
	if next == nil {
		panic("type_error")
	}
	return _partialObject{_sort: nil, _next: next}
}
func (o _partialObject) next() Expr {
	return o._next
}

func (o _partialObject) typeCheck(frame Frame, parent _partialObject) _totalObject {
	// convert a partial _partialObject to a total _partialObject using type-check
	panic("not implemented")
}

type _totalObject struct {
	_partialObject
}

func newTotalObject(sort sorts.Sort, next Expr) _totalObject {
	if sort == nil || next == nil {
		panic("type_error")
	}
	return _totalObject{_partialObject{_sort: sort, _next: next}}
}

func (o _totalObject) parent() _totalObject {
	panic("not implemented")
}
func (o _totalObject) partial() _partialObject {
	return o._partialObject
}

type Frame struct {
	dict ordered_map.OrderedMap[Term, _totalObject]
}

func (frame Frame) set(key Term, o _totalObject) Frame {
	panic("not implemented")
}
func (frame Frame) get(key Term) _totalObject {
	panic("not implemented")
}
