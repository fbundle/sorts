package sorts

type Beta struct {
	Cmd Pi
	Arg Sort
}

func (s Beta) Form() Form {
	return List{s.Cmd.Form(), s.Arg.Form()}
}

func (s Beta) Compile(ctx Context) Sort {
	subCtx := ctx.Set(s.Cmd.Param.Name, s.Arg)
	_ = s.Cmd.Body.Compile(subCtx)
	return s
}
func (s Beta) Level(ctx Context) int {
	panic("not_implemented")
}
func (s Beta) Parent(ctx Context) Sort {
	panic("not_implemented")
}
func (s Beta) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Beta) Reduce(ctx Context) Sort {
	subCtx := ctx.Set(s.Cmd.Param.Name, s.Arg)
	return s.Cmd.Body.Reduce(subCtx)
}
