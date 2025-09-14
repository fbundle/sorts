package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/form"
)

type ListCompileFunc = func(r Context, list form.List) (Context, AlmostSort)

func ListCompileBeta(ctx Context, list form.List) (Context, AlmostSort) {
	if not(len(list) == 2) {
		panic("beta must be (cmd arg)")
	}

	ctx, cmd := ctx.Compile(list[0])
	ctx, arg := ctx.Compile(list[1])

	return ctx, Beta{
		Cmd: cmd,
		Arg: arg,
	}
}

// Beta - beta reduction
type Beta struct {
	Cmd AlmostSort
	Arg AlmostSort
}

func (b Beta) attrAlmostSort(ctx Context) attrAlmostSort {
	return attrAlmostSort{
		form: form.List{Form(ctx, b.Cmd), Form(ctx, b.Arg)},
	}
}

func (b Beta) TypeCheck(ctx Context, parent ActualSort) ActualSort {
	//TODO implement me
	ctx.NewTerm(form.Name(fmt.Sprintf("(%s, %s)")), parent)

	panic("implement me")
}

func ListCompileLambda(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) (Context, AlmostSort) {
		mustMatchHead(Head, list)
		if not(len(list) == 3) {
			panic(fmt.Errorf("lambda must be (%s param body)", Head))
		}
		ctx, body := ctx.Compile(list[0])

		return ctx, Lambda{
			Head:  Head,
			Param: list[0].(form.Name),
			Body:  body,
		}
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Head  form.Name
	Param form.Name
	Body  AlmostSort
}

func (l Lambda) attrAlmostSort(ctx Context) attrAlmostSort {
	return attrAlmostSort{
		form: form.List{l.Head, l.Param, Form(ctx, l.Body)},
	}
}
func (l Lambda) TypeCheck(ctx Context, parent ActualSort) ActualSort {
	//TODO implement me
	panic("implement me")
}

func ListCompileLet(Undef form.Name) func(Head form.Name) ListCompileFunc {
	return func(H form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) (Context, AlmostSort) {
			mustMatchHead(H, list)
			if len(list) < 2 || not((len(list)+1)%3 == 0) {
				panic(fmt.Errorf("let must be (%s name1 type1 value1 ... nameN typeN valueN final)", H))
			}

			bindings := make([]LetBinding, 0)

			var almostType AlmostSort
			for i := 1; i < len(list)-1; i += 3 {
				nameForm, typeForm, valueForm := list[i], list[i+1], list[i+2]

				name, ok := nameForm.(form.Name)
				if !ok {
					panic(TypeErr)
				}

				ctx, almostType = ctx.Compile(typeForm)
				actualType := mustSort(almostType)

				var actualValue ActualSort
				if nameUndef, ok := valueForm.(form.Name); ok && nameUndef == Undef {
					actualValue = ctx.NewTerm(name, actualType)
				} else {
					var almostValue AlmostSort
					ctx, almostValue = ctx.Compile(valueForm)
					actualValue = almostValue.TypeCheck(ctx, actualType)
				}

				b := LetBinding{
					Name:  name,
					Type:  actualType,
					Value: actualValue,
				}

				ctx = ctx.Set(b.Name, b.Value) // memorize for the next parsing

				bindings = append(bindings, b)
			}

			ctx, final := ctx.Compile(list[len(list)-1])

			return ctx, Let{
				Bindings: bindings,
				Final:    final,
			}
		}
	}
}

type LetBinding struct {
	Name  form.Name
	Type  ActualSort
	Value ActualSort
}
type Let struct {
	Head     form.Name
	Bindings []LetBinding
	Final    AlmostSort
}

func (l Let) attrAlmostSort(ctx Context) attrAlmostSort {
	f := form.List{l.Head}
	for _, b := range l.Bindings {
		f = append(f, b.Name)
		f = append(f, Form(ctx, b.Type))
		f = append(f, Form(ctx, b.Value))
	}
	f = append(f, Form(ctx, l.Final))
	return attrAlmostSort{
		form: f,
	}
}
func (l Let) TypeCheck(ctx Context, parent ActualSort) ActualSort {
	panic("implement me")
}

func ListCompileMatch(Exact form.Name) func(H form.Name) ListCompileFunc {
	return func(H form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) (Context, AlmostSort) {
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
					ctx, pattern = ctx.Compile(patternForm[0])
				} else {
					// pattern match
					pattern = patternForm[1]
				}

				var value AlmostSort
				ctx, value = ctx.Compile(valueForm)
				cases = append(cases, MatchCase{
					Pattern: pattern,
					Value:   value,
				})
			}

			var cond AlmostSort
			ctx, cond = ctx.Compile(list[1])
			var final AlmostSort
			ctx, final = ctx.Compile(list[len(list)-1])

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
	Value   AlmostSort
}
type Match struct {
	Head  form.Name
	Cond  AlmostSort
	Cases []MatchCase
	Final AlmostSort
}

func (m Match) attrAlmostSort(ctx Context) attrAlmostSort {}
func (m Match) TypeCheck(ctx Context, parent ActualSort) ActualSort {
	panic("implement me")
}
