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
		if len(list) < 2 {
			panic(fmt.Errorf("inhabitant must be (%s type)", Head))
		}
		parentForm := list[1]
		parent := ctx.Compile(parentForm)

		name := Head + "_" + form.Name(randString(6))
		atom := ctx.NewTerm(name, parent)
		return Inhabitant{
			Atom: atom,
			Head: Head,
			Name: name,
		}
	}
}

type LetBinding struct {
	Name  form.Name
	Value Sort
}

type Let struct {
	Atom
	Head     form.Name
	Bindings []LetBinding
	Final    Sort
}

func ListCompileLet(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) Sort {
		mustMatchHead(Head, list)
		if len(list) < 2 || not(len(list)%2 == 0) {
			panic(fmt.Errorf("let must be (%s name1 value1 ... nameN valueN final)", Head))
		}

		bindings := make([]LetBinding, 0)
		for i := 1; i < len(list)-1; i += 2 {
			nameForm, valueForm := list[i], list[i+1]

			name, ok := nameForm.(form.Name)
			if !ok {
				panic(TypeErr)
			}
			value := ctx.Compile(valueForm)
			bindings = append(bindings, LetBinding{
				Name:  name,
				Value: value,
			})

			ctx = ctx.Set(name, value) // binding
		}

		finalForm := list[len(list)-1]
		final := ctx.Compile(finalForm)
		atom := ctx.NewTerm(list, ctx.Parent(final))

		return Let{
			Atom:     atom,
			Head:     Head,
			Bindings: bindings,
			Final:    final,
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
		return pattern, ctx.Form(mc.Value)
	case Sort: // exact matching
		return ctx.Form(pattern), ctx.Form(mc.Value)
	default:
		panic(TypeErr)
	}
}

type Match struct {
	Atom
	Head  form.Name
	Cond  Sort
	Cases []MatchCase
	Final Sort
}

func ListCompileMatch(Exact form.Name) func(H form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			mustMatchHead(Head, list)
			if len(list) < 3 || not(len(list)%2 == 1) {
				panic(fmt.Errorf("match must be (%s cond pattern1 value1 ... patternN valueN final)", Head))
			}
			condForm := list[1]
			cond := ctx.Compile(condForm)

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
					pattern = ctx.Compile(patternForm[1])
					value = ctx.Compile(valueForm)
				} else {
					// pattern match
					pattern = patternForm[1]
					// suppose the pattern is (succ z) - how to set z into nextCtx to compile valueForm?
					panic("TODO")
				}

				cases = append(cases, MatchCase{
					Pattern: pattern,
					Value:   value,
				})
			}

			finalForm := list[len(list)-1]
			final := ctx.Compile(finalForm)

			// find the weakest type / largest type of return values
			typeList := []Sort{ctx.Parent(final)}
			for _, c := range cases {
				typeList = append(typeList, ctx.Parent(c.Value))
			}

			returnType := func() Sort {
				for _, t1 := range typeList {
					dominate := true
					// check if c1Type dominate everyone else
					for _, t2 := range typeList {
						if not(ctx.LessEqual(t2, t1)) {
							dominate = false
							break
						}
					}
					if dominate {
						return t1
					}
				}
				// return terminal type of maximal level
				maxLevel := ctx.Level(typeList[0])
				for _, t := range typeList {
					maxLevel = max(maxLevel, ctx.Level(t))
				}
				return ctx.Terminal(maxLevel)
			}()

			atom := ctx.NewTerm(list, returnType)

			return Match{
				Atom:  atom,
				Head:  Head,
				Cond:  cond,
				Cases: cases,
				Final: final,
			}
		}
	}
}
