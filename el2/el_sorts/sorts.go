package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type ListCompileFunc = func(r Context, list form.List) Sort

// Beta - beta reduction
type Beta struct {
	Atom
	Cmd Sort
	Arg Sort
}

func ListCompileBeta(ctx Context, list form.List) (Context, Sort) {
	if not(len(list) == 2) {
		panic("beta must be (cmd arg)")
	}

	cmd := ctx.Compile(list[0])
	arg := ctx.Compile(list[1])

	// type check
	arrow, ok := ctx.Parent(cmd).(sorts.Arrow)
	if !ok {
		panic(TypeErr)
	}
	mustTermOf(ctx, arg, arrow.A)
	atom := ctx.NewTerm(list, arrow.B)

	return ctx, Beta{
		Atom: atom,
		Cmd:  cmd,
		Arg:  arg,
	}
}

type LambdaParam struct {
	Name form.Name
	Type Sort
}

// Lambda - lambda abstraction
type Lambda struct {
	Atom
	Head   form.Name
	Params []LambdaParam
	Body   Sort
}

func ListCompileLambda(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) Sort {
		mustMatchHead(Head, list)
		if len(list) < 2 {
			panic(fmt.Errorf("lambda must be (%s param1 type1 ... paramN typeN body)", Head))
		}

		params := make([]LambdaParam, 0)
		for i := 2; i < len(list)-1; i += 2 {
			paramForm, paramTypeForm := list[i], list[i+1]
			param, ok := paramForm.(form.Name)
			if !ok {
				panic(TypeErr)
			}
			params = append(params, LambdaParam{
				Name: param,
				Type: ctx.Compile(paramTypeForm),
			})
		}

		// suppose params is of the correct type, compile body
		for _, param := range params {
			ctx = ctx.Set(param.Name, ctx.NewTerm(param.Name, param.Type))
		}

		bodyForm := list[len(list)-1]
		body := ctx.Compile(bodyForm)

		arrow := ctx.Parent(body)
		for i := len(params) - 1; i >= 0; i-- {
			arrow = sorts.Arrow{
				A: params[i].Type,
				B: arrow,
			}
		}
		atom := ctx.NewTerm(list, arrow)

		return Lambda{
			Atom:   atom,
			Head:   Head,
			Params: params,
			Body:   body,
		}

	}
}

// Inhabitant - give an undefined term of a type
type Inhabitant struct {
	Atom
	Head form.Name
	Name form.Name
}

func ListCompileInhabitant(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) Sort {
		mustMatchHead(Head, list)
		if len(list) < 3 {
			panic(fmt.Errorf("inhabitant must be (%s name type)", Head))
		}
		nameForm := list[1]
		name, ok := nameForm.(form.Name)
		if !ok {
			panic(TypeErr)
		}
		parentForm := list[2]
		parent := ctx.Compile(parentForm)

		atom := ctx.NewTerm(name, parent)
		return Inhabitant{
			Atom: atom,
			Head: Head,
			Name: name,
		}
	}
}

type MatchCase struct {
	Pattern any // Union[form.Form, Sort] - pattern matching vs exact matching
	Value   Sort
}

func (mc MatchCase) form(ctx Context) (form.Form, form.Form) {
	switch pattern := mc.Pattern.(type) {
	case form.Form: // pattern matching
		return pattern, Form(ctx, mc.Value)
	case Sort: // exact matching
		return Form(ctx, pattern), Form(ctx, mc.Value)
	default:
		panic(TypeErr)
	}
}

type Match struct {
	Head  form.Name
	Cond  Sort
	Cases []MatchCase
	Final Sort
}

func ListCompileMatch(Exact form.Name) func(H form.Name) ListCompileFunc {
	return func(H form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			mustMatchHead(H, list)
			if len(list) < 3 || not(len(list)%2 == 1) {
				panic(fmt.Errorf("match must be (%s cond pattern1 value1 ... patternN valueN final)", H))
			}
			var cond Sort
			ctx, cond = ctx.Compile(list[1])

			cases := make([]MatchCase, 0)
			for i := 2; i < len(list)-1; i += 2 {
				patternForm, valueForm := list[i], list[i+1]

				var pattern any
				var value Sort
				if patternForm, ok := patternForm.(form.List); ok && patternForm[0] == Exact {
					if len(patternForm) != 2 {
						panic(fmt.Errorf("exact match must be (%s form)", Exact))
					}
					// exact match
					ctx, pattern = ctx.Compile(patternForm[1])
					ctx, value = ctx.Compile(valueForm)
				} else {
					// pattern match
					pattern = patternForm[1]
					// suppose the pattern is (succ z) - how to set z into nextCtx to compile valueForm?
					panic("checkpoint")
				}

				cases = append(cases, MatchCase{
					Pattern: pattern,
					Value:   value,
				})
			}

			var final Sort
			ctx, final = ctx.Compile(list[len(list)-1])

			return ctx, Match{
				Cond:  cond,
				Cases: cases,
				Final: final,
			}
		}
	}
}
