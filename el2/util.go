package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Name = form.Name
type Form = form.Form
type List = form.List
type Sort = sorts.Sort
type SortAttr = sorts.SortAttr

var NewAtomChain = sorts.NewAtomChain
var NewAtomTerm = sorts.NewAtomTerm

type ParseFunc = func(form Form) almost_sort.AlmostSort
type ListParseFunc = func(parse ParseFunc, list List) almost_sort.AlmostSort

var ListParseArrow = sorts.ListParseArrow
var ListParseSum = sorts.ListParseSum
var ListParseProd = sorts.ListParseProd

var TypeErr = fmt.Errorf("type_error")

func toListParser(listParse sorts.ListParseFunc) ListParseFunc {
	return func(parse ParseFunc, list List) almost_sort.AlmostSort {
		sort := listParse(func(form sorts.Form) sorts.Sort {
			s := almost_sort.mustSort(parse(form)) // inside a sort, must be sort
			if s == nil {
				panic(TypeErr)
			}
			return s
		}, list)
		return almost_sort.ActualSort{sort}
	}
}
