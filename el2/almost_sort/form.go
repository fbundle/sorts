package el_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func ListParseBeta(parse ParseFunc, list form.List) AlmostSort {
	if not(len(list) == 2) {
		panic("beta must be (cmd arg)")
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

func (b Beta) almostSortAttr() {}

func (b Beta) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	//TODO implement me
	panic("implement me")
}

func ListParseLambda(H form.Name) ListParseFunc {
	return func(parse ParseFunc, list form.List) AlmostSort {
		mustMatchHead(H, list)
		if not(len(list) == 3) {
			panic(fmt.Errorf("lambda must be (%s param body)", H))
		}
		return Lambda{
			Param: list[0].(form.Name),
			Body:  parse(list[1]),
		}
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Param form.Name
	Body  AlmostSort
}

func (l Lambda) almostSortAttr() {}
func (l Lambda) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	//TODO implement me
	panic("implement me")
}

func ListParseLet(H form.Name) ListParseFunc {
	return func(parse ParseFunc, list form.List) AlmostSort {
		mustMatchHead(H, list)
		if len(list) < 2 || not((len(list)+1)%3 == 0) {
			panic(fmt.Errorf("let must be (%s name1 type1 value1 ... nameN typeN valueN final)", H))
		}

		bindings := make([]LetBinding, 0)
		for i := 1; i < len(list)-1; i += 3 {
			nameForm, typeForm, valueForm := list[i], list[i+1], list[i+2]

			b := LetBinding{
				Name:  nameForm.(form.Name),
				Type:  MustSort(parse(typeForm)),
				Value: parse(valueForm),
			}
			bindings = append(bindings, b)
		}
		return Let{
			Bindings: bindings,
			Final:    parse(list[len(list)-1]),
		}
	}
}

type LetBinding struct {
	Name  form.Name
	Type  sorts.Sort
	Value AlmostSort
}
type Let struct {
	Bindings []LetBinding
	Final    AlmostSort
}

func (l Let) almostSortAttr() {}
func (l Let) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	panic("implement me")
}

func ListParseMatch(H form.Name) ListParseFunc {
	return func(parse ParseFunc, list form.List) AlmostSort {
		mustMatchHead(H, list)
		if len(list) < 2 || not(len(list)%2 == 0) {
			panic(fmt.Errorf("match must be (%s pattern1 value1 ... patternN valueN final)", H))
		}

		return Match{}
	}
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

func (m Match) almostSortAttr() {}
func (m Match) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	panic("implement me")
}
