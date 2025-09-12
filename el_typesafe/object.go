package el_typesafe

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type object struct {
	sort sorts.Sort
	next Expr
}

type totalObject struct {
	object
}

func newTotalObject(sort sorts.Sort, next Expr) totalObject {
	if sort == nil || next == nil {
		panic("type_error")
	}
	return totalObject{object{sort: sort, next: next}}
}

func (o totalObject) parent() totalObject {
	panic("not implemented")
}
func (o totalObject) partial() partialObject {
	return partialObject{object: o.object}
}

type partialObject struct {
	object
}

func newPartialObject(next Expr) partialObject {
	if next == nil {
		panic("type_error")
	}
	return partialObject{object{sort: nil, next: next}}
}

func (o partialObject) typeCheck(frame Frame, parentSort sorts.Sort) totalObject {
	// convert a partial object to a total object using type-check
	panic("not implemented")
}

type Frame struct {
	dict ordered_map.OrderedMap[Term, totalObject]
}

func (frame Frame) set(key Term, o totalObject) {
	panic("not implemented")
}
func (frame Frame) get(key Term) (totalObject, bool) {
	panic("not implemented")
}
