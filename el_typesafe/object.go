package el_typesafe

import (
	"errors"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

var typeErr = errors.New("type_error")

type _partialObject struct {
	_sort sorts.Sort
	next  Expr
}

func newPartialObject(next Expr) _partialObject {
	if next == nil {
		panic(typeErr)
	}
	return _partialObject{_sort: nil, next: next}
}

func (o _partialObject) typeCheck(frame Frame, parent _totalObject) _totalObject {
	// convert a partial _partialObject to a total _partialObject using type-check
	panic("not implemented")
}

type _totalObject struct {
	_sort sorts.Sort
	next  totalExpr
}

func newTotalObject(sort sorts.Sort, next totalExpr) _totalObject {
	if sort == nil || next == nil {
		panic(typeErr)
	}
	return _totalObject{_sort: sort, next: next}
}

func (o _totalObject) parent() _totalObject {
	panic("not implemented")
}
func (o _totalObject) partial() _partialObject {
	panic("not implemented")
}

type Frame struct {
	dict ordered_map.OrderedMap[Term, _totalObject]
}

func (frame Frame) set(key Term, o _totalObject) Frame {
	return Frame{dict: frame.dict.Set(key, o)}
}
func (frame Frame) get(key Term) _totalObject {
	o, ok := frame.dict.Get(key)
	if !ok {
		panic(typeErr)
	}
	return o
}
