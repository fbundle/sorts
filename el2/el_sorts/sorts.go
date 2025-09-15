package el_sorts

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

func ListCompileBeta(ctx Context, list form.List) Sort {
	if not(len(list) == 2) {
		panic(fmt.Errorf("beta must be (cmd arg)"))
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

	return Beta{
		Atom: atom,
		Cmd:  cmd,
		Arg:  arg,
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Atom
	Head  form.Name
	Param TypeAnnot
	Body  Sort
}

func ListCompileLambda(TypeAnnot form.Name) func(Head form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			err := fmt.Errorf("lambda must be (%s (%s param type) body)", Head, TypeAnnot)
			mustMatchHead(err, Head, list)
			if len(list) != 3 {
				panic(err)
			}

			param := ListParseTypeAnnot(TypeAnnot)(ctx, mustList(err, list[1]))

			ctx = ctx.Set(param.Name, ctx.NewTerm(param.Name, param.Type))
			bodyForm := list[len(list)-1]
			body := ctx.Compile(bodyForm)

			atom := ctx.NewTerm(list, sorts.Arrow{
				A: param.Type,
				B: ctx.Parent(body),
			})

			return Lambda{
				Atom:  atom,
				Head:  Head,
				Param: param,
				Body:  body,
			}

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
		err := fmt.Errorf("inhabitant must be (%s type)", Head)
		mustMatchHead(err, Head, list)
		if len(list) < 2 {
			panic(err)
		}
		parentForm := list[1]
		parent := ctx.Compile(parentForm)

		name := form.Name(fmt.Sprintf("%s_%d", Head, nextValue()))
		atom := ctx.NewTerm(name, parent)
		return Inhabitant{
			Atom: atom,
			Head: Head,
			Name: name,
		}
	}
}

type Let struct {
	Atom
	Head     form.Name
	Bindings []NameBinding
	Final    Sort
}

func ListCompileLet(Assign form.Name) func(Head form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			err := fmt.Errorf("let must be (%s (%s name1 value1) ... (%s nameN valueN) final)", Head, Assign, Assign)
			mustMatchHead(err, Head, list)
			if len(list) < 2 || not(len(list)%2 == 0) {
				panic(err)
			}

			bindings := make([]NameBinding, 0)
			for i := 1; i < len(list)-1; i++ {
				binding := ParseNameBinding(Assign)(ctx, mustList(err, list[i]))
				bindings = append(bindings, binding)

				// binding
				ctx = ctx.Set(binding.Name, binding.Value)
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
}

type Match struct {
	Atom
	Head  form.Name
	Cond  Sort
	Cases []MathCase
	Final Sort
}

func ListCompileMatch(Exact form.Name) func(H form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			err := fmt.Errorf("match must be (%s cond (pattern1 value1) ... patternN valueN final)", Head)

			mustMatchHead(Head, list)
			if len(list) < 3 || not(len(list)%2 == 1) {
				panic(fmt.Errorf("match must be (%s cond pattern1 value1 ... patternN valueN final)", Head))
			}
			condForm := list[1]
			cond := ctx.Compile(condForm)

			cases := make([]MathCase, 0)
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
					panic(fmt.Errorf("pattern matching is not implemented for now"))
				}

				cases = append(cases, MathCase{
					Pattern: pattern,
					Value:   value,
				})
			}

			finalForm := list[len(list)-1]
			final := ctx.Compile(finalForm)

			// find the weakest type / largest type of return values
			typeList := make([]Sort, 0, len(cases)+1)
			for _, c := range cases {
				typeList = append(typeList, ctx.Parent(c.Value))
			}
			typeList = append(typeList, ctx.Parent(final))

			returnType := sorts.LeastUpperBound(ctx, "⊕", typeList...)

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
