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
	arrowType := mustType[Pi](TypeError, s.Cmd.Parent(ctx))
	argType := s.Arg.Parent(ctx)
	if !argType.LessEqual(ctx, arrowType.Param.Type) {
		panic(TypeError)
	}

	return Beta{
		Cmd: arrowType,
		Arg: s.Arg,
	}
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
