package sorts

type Beta struct {
	Cmd Pi
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
func (s Beta) Reduce(ctx Context) Sort {
	return s.Cmd.Body.Reduce(ctx.Set(s.Cmd.Param.Name, s.Arg))
}
