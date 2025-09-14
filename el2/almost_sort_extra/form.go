package el2_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func ListParseBeta(ctx Runtime, list form.List) (Runtime, almost_sort.AlmostSort) {
	if not(len(list) == 2) {
		panic("beta must be (cmd arg)")
	}

	ctx, cmd := ctx.Parse(list[0])
	ctx, arg := ctx.Parse(list[1])

	return ctx, Beta{
		Cmd: cmd,
		Arg: arg,
	}
}

// Beta - beta reduction
type Beta struct {
	Cmd almost_sort.AlmostSort
	Arg almost_sort.AlmostSort
}

func (b Beta) almostSortAttr() {}

func (b Beta) TypeCheck(sa sorts.SortAttr, parent almost_sort.ActualSort) almost_sort.ActualSort {
	//TODO implement me
	panic("implement me")
}

func ListParseLambda(H form.Name) ListParseFunc {
	return func(ctx Runtime, list form.List) (Runtime, almost_sort.AlmostSort) {
		mustMatchHead(H, list)
		if not(len(list) == 3) {
			panic(fmt.Errorf("lambda must be (%s param body)", H))
		}
		ctx, body := ctx.Parse(list[0])

		return ctx, Lambda{
			Param: list[0].(form.Name),
			Body:  body,
		}
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Param form.Name
	Body  almost_sort.AlmostSort
}

func (l Lambda) almostSortAttr() {}
func (l Lambda) TypeCheck(sa sorts.SortAttr, parent almost_sort.ActualSort) almost_sort.ActualSort {
	//TODO implement me
	panic("implement me")
}

func ListParseLet(Undef form.Name) func(H form.Name) ListParseFunc {
	return func(H form.Name) ListParseFunc {
		return func(ctx Runtime, list form.List) (Runtime, almost_sort.AlmostSort) {
			mustMatchHead(H, list)
			if len(list) < 2 || not((len(list)+1)%3 == 0) {
				panic(fmt.Errorf("let must be (%s name1 type1 value1 ... nameN typeN valueN final)", H))
			}

			bindings := make([]LetBinding, 0)

			var almostType almost_sort.AlmostSort
			for i := 1; i < len(list)-1; i += 3 {
				nameForm, typeForm, valueForm := list[i], list[i+1], list[i+2]

				var almostValue almost_sort.AlmostSort
				if valueForm, ok := valueForm.(form.Name); ok && valueForm == Undef {
					almostValue = nil
				} else {
					ctx, almostValue = ctx.Parse(valueForm)
				}

				ctx, almostType = ctx.Parse(typeForm)

				actualType := almost_sort.MustSort(almostType)
				actualValue := almostValue.TypeCheck(ctx.SortAttr(), actualType)

				b := LetBinding{
					Name:  nameForm.(form.Name),
					Type:  actualType,
					Value: actualValue,
				}

				ctx = ctx.Set(b.Name, b.Value) // memorize for the next parsing

				bindings = append(bindings, b)
			}

			ctx, final := ctx.Parse(list[len(list)-1])

			return ctx, Let{
				Bindings: bindings,
				Final:    final,
			}
		}
	}
}

type LetBinding struct {
	Name  form.Name
	Type  almost_sort.ActualSort
	Value almost_sort.ActualSort
}
type Let struct {
	Bindings []LetBinding
	Final    almost_sort.AlmostSort
}

func (l Let) almostSortAttr() {}
func (l Let) TypeCheck(sa sorts.SortAttr, parent almost_sort.ActualSort) almost_sort.ActualSort {
	panic("implement me")
}

func ListParseMatch(Exact form.Name) func(H form.Name) ListParseFunc {
	return func(H form.Name) ListParseFunc {
		return func(ctx Runtime, list form.List) (Runtime, almost_sort.AlmostSort) {
			mustMatchHead(H, list)
			if len(list) < 3 || not(len(list)%2 == 1) {
				panic(fmt.Errorf("match must be (%s cond pattern1 value1 ... patternN valueN final)", H))
			}

			cases := make([]MatchCase, 0)
			for i := 2; i < len(list)-1; i += 2 {
				patternForm, valueForm := list[i], list[i+1]

				var pattern any
				if patternForm, ok := patternForm.(form.List); ok && patternForm[0] == Exact {
					if len(patternForm) != 2 {
						panic(fmt.Errorf("exact match must be (%s form)", Exact))
					}
					// exact match
					ctx, pattern = ctx.Parse(patternForm[0])
				} else {
					// pattern match
					pattern = patternForm[1]
				}

				var value almost_sort.AlmostSort
				ctx, value = ctx.Parse(valueForm)
				cases = append(cases, MatchCase{
					Pattern: pattern,
					Value:   value,
				})
			}

			var cond almost_sort.AlmostSort
			ctx, cond = ctx.Parse(list[1])
			var final almost_sort.AlmostSort
			ctx, final = ctx.Parse(list[len(list)-1])

			return ctx, Match{
				Cond:  cond,
				Cases: cases,
				Final: final,
			}
		}
	}
}

type MatchCase struct {
	Pattern any // Union[form.Form, ActualSort] - pattern matching vs exact matching
	Value   almost_sort.AlmostSort
}
type Match struct {
	Cond  almost_sort.AlmostSort
	Cases []MatchCase
	Final almost_sort.AlmostSort
}

func (m Match) almostSortAttr() {}
func (m Match) TypeCheck(sa sorts.SortAttr, parent almost_sort.ActualSort) almost_sort.ActualSort {
	panic("implement me")
}
