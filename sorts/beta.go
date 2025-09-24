package sorts

func init() {
	DefaultParseFunc = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		return Beta{
			Cmd: mustType[Lambda](err, ctx.Parse(list[0])),
			Arg: ctx.Parse(list[1]),
		}
	}
}

type Beta struct {
	Cmd Lambda
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
