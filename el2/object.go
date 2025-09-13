package el2

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Name = form.Name
type Form = form.Form
type List = form.List
type Sort = sorts.Sort

var NewAtomChain = sorts.NewAtomChain
var NewAtomTerm = sorts.NewAtomTerm

type ListParseFunc = sorts.ListParseFunc

var ListParseArrow = sorts.ListParseArrow
var ListParseSum = sorts.ListParseSum
var ListParseProd = sorts.ListParseProd

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	Sort
	TypeCheck(parent Sort) Sort
}
