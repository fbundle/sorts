package sorts

import "github.com/fbundle/sorts/slices_util"

func init() {
	DefaultParseFunc = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		output := ctx.Parse(list[0])
		args := slices_util.Map(list[1:], func(form Form) Sort {
			return ctx.Parse(form)
		})
		slices_util.ForEach(args, func(arg Sort) {
			output = Beta{
				Cmd: output,
				Arg: arg,
			}
		})
		return output
	}
}

type Beta struct {
	Cmd Sort
	Arg Sort
}

func (s Beta) Form() Form {
	return List{s.Cmd.Form(), s.Arg.Form()}
}

func (s Beta) Parent(ctx Context) Sort {
	return s.Cmd.Body.Parent(ctx.Set(s.Cmd.Param.Name, s.Arg))
}
func (s Beta) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Beta) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Beta) Eval(ctx Context) Sort {
	return s.Cmd.Body.Eval(ctx.Set(s.Cmd.Param.Name, s.Arg))
}

func init() {
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(LambdaCmd), "param1", "...", "paramN", "body"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		params := slices_util.Map(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(ctx, mustType[List](err, form))
		})
		output := ctx.Parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(params), func(param Annot) {
			output = Lambda{
				Param: param,
				Body:  output,
			}
		})
		return output
	}
	const ArrowCmd Name = "->"
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) Sort {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		// TODO - probably we don't need this anymore
		panic("not implemented")
	}
}

const (
	LambdaCmd Name = "=>"
)

// Lambda - lambda abstraction (or Pi-type)
type Lambda struct {
	Param Annot
	Body  Sort
}

func (s Lambda) Form() Form {
	return List{LambdaCmd, s.Param.Form(), s.Body.Form()}
}

func (s Lambda) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Lambda) Parent(ctx Context) Sort {
	return Lambda{
		Param: s.Param,
		Body: Type{
			Body: s.Body,
		},
	}
}
func (s Lambda) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Lambda) Eval(ctx Context) Sort {
	return s
}
