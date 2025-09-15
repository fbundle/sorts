package sorts

import (
	"fmt"

	"github.com/fbundle/sorts/form"
)

var TypeCheckErr = fmt.Errorf("type_check")

type ReducibleSort interface {
	Sort
	Reduce(ctx Context) Sort
}

// Beta - beta reduction
type Beta struct {
	Atom
	Cmd Sort
	Arg Sort
}

func (b Beta) Reduce(ctx Context) Sort {
	if ctx.Mode() != ModeEval {
		return b
	}

	// TODO - implement evaluation
	return b
}

func ListCompileBeta(ctx Context, list form.List) Sort {
	err := fmt.Errorf("beta must be (cmd arg)")
	if len(list) != 2 {
		panic(err)
	}

	cmd := ctx.Compile(list[0])
	arg := ctx.Compile(list[1])

	// type check
	arrow, ok := ctx.Parent(cmd).(Arrow)
	if !ok {
		panic(err)
	}

	if !ctx.LessEqual(ctx.Parent(arg), arrow.A) {
		panic(TypeCheckErr)
	}

	atom := ctx.NewTerm(list, arrow.B)

	return (Beta{
		Atom: atom,
		Cmd:  cmd,
		Arg:  arg,
	}).Reduce(ctx)
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

			atom := ctx.NewTerm(list, Arrow{
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

		name := form.Name(fmt.Sprintf("%s_%d", Head, nextCount()))
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

func (l Let) Reduce(ctx Context) Sort {
	if ctx.Mode() != ModeEval {
		return l
	}

	// TODO - implement evaluation
	return l
}

func ListCompileLet(Assign form.Name) func(Head form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			err := fmt.Errorf("let must be (%s (%s name1 value1) ... (%s nameN valueN) final)", Head, Assign, Assign)
			mustMatchHead(err, Head, list)
			if len(list) < 2 {
				panic(err)
			}

			bindings := make([]NameBinding, 0)
			for i := 1; i < len(list)-1; i++ {
				binding := ParseNameBinding(Assign)(
					ctx, mustList(err, list[i]),
				)
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
	Cases []MatchCase
}

func (m Match) Reduce(ctx Context) Sort {
	if ctx.Mode() != ModeEval {
		return m
	}

	// TODO - implement evaluation
	return m
}

func ListCompileMatch(Arrow form.Name, DefaultConstr form.Name) func(Head form.Name) ListCompileFunc {
	return func(Head form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) Sort {
			err := fmt.Errorf("match must be (%s cond (%s pattern1 value1) ... (%s patternN valueN))", Head, Arrow, Arrow)
			mustMatchHead(err, Head, list)

			if len(list) < 2 {
				panic(err)
			}
			condForm := list[1]
			cond := ctx.Compile(condForm)

			cases := make([]MatchCase, 0)
			for i := 2; i < len(list); i++ {
				cases = append(cases, ParseMatchCase(Arrow, DefaultConstr, ctx.Parent(cond))(
					ctx, mustList(err, list[i]),
				))
			}

			// find the weakest type / largest type of return values
			typeList := make([]Sort, 0, len(cases))
			for _, c := range cases {
				typeList = append(typeList, ctx.Parent(c.Value))
			}

			returnType := LeastUpperBound(ctx, "⊕", typeList...)

			atom := ctx.NewTerm(list, returnType)

			return Match{
				Atom:  atom,
				Head:  Head,
				Cond:  cond,
				Cases: cases,
			}
		}
	}
}

type Inspect struct {
	Atom
	Head  form.Name
	Value Sort
}

func ListCompileInspect(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) Sort {
		err := fmt.Errorf("inspect must be (%s value)", Head)
		mustMatchHead(err, Head, list)
		if len(list) != 2 {
			panic(err)
		}

		value := ctx.Compile(list[1])

		// do inspect
		fmt.Println("inspect", ctx.ToString(value))

		atom := ctx.NewTerm(list, ctx.Parent(value))
		return Inspect{
			Atom:  atom,
			Head:  Head,
			Value: value,
		}

	}
}
