package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

func NewActualSort(sort sorts.Sort) ActualSort {
	return ActualSort{_sort: sort}
}

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	attrAlmostSort(ctx Context) attrAlmostSort
	TypeCheck(ctx Context, parent ActualSort) ActualSort
}

type attrAlmostSort struct {
	form form.Form
}

// ActualSort - a sort
type ActualSort struct {
	_sort sorts.Sort
}

func (s ActualSort) attrAlmostSort(ctx Context) attrAlmostSort {
	return attrAlmostSort{
		form: ctx.Form(s._sort),
	}
}

func (s ActualSort) Repr() sorts.Sort {
	return s._sort
}

func (s ActualSort) TypeCheck(ctx Context, parent ActualSort) ActualSort {
	if !ctx.LessEqual(ctx.Parent(s._sort), parent._sort) {
		panic(TypeErr)
	}
	return ActualSort{s._sort}
}
