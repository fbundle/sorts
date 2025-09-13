package el2

import (
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

type ListParseFunc = sorts.ListParseFunc

var ListParseArrow = sorts.ListParseArrow
var ListParseSum = sorts.ListParseSum
var ListParseProd = sorts.ListParseProd

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	TypeCheck(sa SortAttr, parent Sort) Sort
}

type actualSort struct {
	Sort
}

func (s actualSort) TypeCheck(sa SortAttr, parent Sort) Sort {
	must(sa).termOf(s, parent)
	return s.Sort
}
