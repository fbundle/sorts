package el_almost_sort

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func ListParseBeta(parse ParseFunc, list form.List) AlmostSort {
	if len(list) != 2 {
		panic("beta list must have two elements")
	}
	return Beta{
		Cmd: parse(list[0]),
		Arg: parse(list[1]),
	}
}

// Beta - beta reduction
type Beta struct {
	Cmd AlmostSort
	Arg AlmostSort
}

func (f Beta) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	//TODO implement me
	panic("implement me")
}

func ListParseLambda(parse ParseFunc, list form.List) AlmostSort {
	if len(list) != 2 {
		panic("lambda list must have two elements")
	}
	return Lambda{
		Param: list[0].(form.Name),
		Body:  parse(list[1]),
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Param form.Name
	Body  AlmostSort
}

func (l Lambda) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	//TODO implement me
	panic("implement me")
}

type LetBinding struct {
	Name  form.Name
	Value AlmostSort
}
type Let struct {
	Bindings []LetBinding
	Final    AlmostSort
}

func (l Let) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	panic("implement me")
}

type MatchCase struct {
	Pattern form.Form
	Value   AlmostSort
}
type Match struct {
	Cond  AlmostSort
	Cases []MatchCase
	Final AlmostSort
}

func (m Match) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	panic("implement me")
}
