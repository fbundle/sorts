package el_typesafe

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type _object struct {
	_sort sorts.Sort
	_next Expr
}

func (o _object) next() Expr {
	return o._next
}

type _totalObject struct {
	_object
}

func newTotalObject(sort sorts.Sort, next Expr) _totalObject {
	if sort == nil || next == nil {
		panic("type_error")
	}
	return _totalObject{_object{_sort: sort, _next: next}}
}

func (o _totalObject) parent() _totalObject {
	panic("not implemented")
}
func (o _totalObject) partial() partialObject {
	return partialObject{_object: o._object}
}

type partialObject struct {
	_object
}

func newPartialObject(next Expr) partialObject {
	if next == nil {
		panic("type_error")
	}
	return partialObject{_object{_sort: nil, _next: next}}
}

func (o partialObject) typeCheck(frame Frame, parentSort sorts.Sort) _totalObject {
	// convert a partial _object to a total _object using type-check
	panic("not implemented")
}

type Frame struct {
	dict ordered_map.OrderedMap[Term, _totalObject]
}

func (frame Frame) set(key Term, sort sorts.Sort, next Expr) Frame {
	panic("not implemented")
}
func (frame Frame) get(key Term) _totalObject {
	panic("not implemented")
}
