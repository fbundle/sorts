package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func mustTermOf(ctx Context, x sorts.Sort, X sorts.Sort) {
	if !ctx.LessEqual(ctx.Parent(x), X) {
		panic(TypeErr)
	}
}

type ListCompileFunc = func(r Context, list form.List) (Context, Sort)

func ListCompileBeta(ctx Context, list form.List) (Context, Sort) {
	if not(len(list) == 2) {
		panic("beta must be (cmd arg)")
	}

	ctx, cmd := ctx.Compile(list[0])
	ctx, arg := ctx.Compile(list[1])

	// type check
	arrow, ok := ctx.Parent(cmd.Sort()).(sorts.Arrow)
	if !ok {
		panic(TypeErr)
	}
	mustTermOf(ctx, arg.Sort(), arrow.A)
	atom :=

	return ctx, Beta{

		Cmd: cmd,
		Arg: arg,
	}
}

// Beta - beta reduction
type Beta struct {
	Atom sorts.Atom
	Cmd  Sort
	Arg  Sort
}

func (b Beta) Sort() sorts.Sort {
	return b.Atom
}

func (b Beta) TypeCheck(ctx Context, parent typeSort) typeSort {
	// TODO implement me
	// type check all pass for now
	return ctx.NewTerm(Form(ctx, b), parent)
}

func ListCompileLambda(Head form.Name) ListCompileFunc {
	return func(ctx Context, list form.List) (Context, Sort) {
		mustMatchHead(Head, list)
		if len(list) < 3 {
			panic(fmt.Errorf("lambda must be (%s type param1 type1 ... paramN typeN body)", Head))
		}
		paramForm, bodyForm := list[1], list[2]
		param := paramForm.(form.Name)

		//
		panic("checkpoint")
		// bodyCtx := ctx.Set(param, ctx.NewTerm())
		ctx, body := ctx.Compile(bodyForm)

		return ctx, Lambda{
			Head:  Head,
			Param: param,
			Body:  body,
		}
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Head   form.Name
	Type   Sort
	Params []form.Name
	Body   Sort
}

func (l Lambda) attrSort(ctx Context) attrAlmostSort {
	return attrAlmostSort{
		form: form.List{l.Head, l.Param, Form(ctx, l.Body)},
	}
}
func (l Lambda) TypeCheck(ctx Context, parent typeSort) typeSort {
	// TODO implement me
	// type check all pass for now
	return ctx.NewTerm(Form(ctx, l), parent)
}

func ListCompileLet(Undef form.Name) func(Head form.Name) ListCompileFunc {
	return func(H form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) (Context, Sort) {
			mustMatchHead(H, list)
			if len(list) < 2 || not((len(list)+1)%3 == 0) {
				panic(fmt.Errorf("let must be (%s name1 type1 value1 ... nameN typeN valueN final)", H))
			}

			bindings := make([]LetBinding, 0)

			var almostType Sort
			for i := 1; i < len(list)-1; i += 3 {
				nameForm, typeForm, valueForm := list[i], list[i+1], list[i+2]

				name, ok := nameForm.(form.Name)
				if !ok {
					panic(TypeErr)
				}

				ctx, almostType = ctx.Compile(typeForm)
				actualType := mustSort(almostType)

				undefValue := ctx.NewTerm(name, actualType)

				var actualValue typeSort
				if nameUndef, ok := valueForm.(form.Name); ok && nameUndef == Undef {
					actualValue = undefValue
				} else {
					// temporary add name with the correct type for recursive function
					// i.e. assuming the name is already type-checked
					recCtx := ctx.Set(name, undefValue)
					// compile value
					recCtx, almostValue := recCtx.Compile(valueForm)
					actualValue = almostValue.TypeCheck(recCtx, actualType)
					// remove name
					recCtx = recCtx.Del(name)
					ctx = recCtx
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
	Type  typeSort
	Value typeSort
}
type Let struct {
	Head     form.Name
	Bindings []LetBinding
	Final    Sort
}

func (l Let) attrSort(ctx Context) attrAlmostSort {
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
func (l Let) TypeCheck(ctx Context, parent typeSort) typeSort {
	// TODO implement me
	// type check all pass for now
	return ctx.NewTerm(Form(ctx, l), parent)
}

func ListCompileMatch(Exact form.Name) func(H form.Name) ListCompileFunc {
	return func(H form.Name) ListCompileFunc {
		return func(ctx Context, list form.List) (Context, Sort) {
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

func (m Match) attrSort(ctx Context) attrAlmostSort {
	f := form.List{m.Head, Form(ctx, m.Cond)}
	for _, c := range m.Cases {
		patternForm, valueForm := c.form(ctx)
		f = append(f, patternForm)
		f = append(f, valueForm)
	}
	f = append(f, Form(ctx, m.Final))

	return attrAlmostSort{
		form: f,
	}
}
func (m Match) TypeCheck(ctx Context, parent typeSort) typeSort {
	// TODO implement me
	// type check all pass for now
	return ctx.NewTerm(Form(ctx, m), parent)
}
