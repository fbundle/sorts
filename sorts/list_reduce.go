package sorts

const (
	BetaCmd Name = "β"
)

func init() {
	ListParseFuncMap[BetaCmd] = func(ctx Context, list List) (Context, Sort) {
		err := parseErr(BetaCmd, []string{"cmd", "arg1", "...", "argN"}, "where N >= 1")
		if len(list) < 2 {
			panic(err)
		}

		ctx, cmd := ctx.Parse(list[0])
		args := make([]Sort, 0, len(list)-1)
		for i := 1; i < len(list); i++ {
			var arg Sort
			ctx, arg = ctx.Parse(list[i])
			args = append(args, arg)
		}

		output := Beta{
			Cmd: cmd,
			Arg: args[0],
		}
		for i := 1; i < len(args); i++ {
			output = Beta{
				Cmd: output,
				Arg: args[i],
			}
		}

		return ctx, output
	}
}

type Beta struct {
	Cmd Sort
	Arg Sort
}

func (s Beta) Form() Form {
	return List{BetaCmd, s.Cmd.Form(), s.Arg.Form()}
}

func (s Beta) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Beta) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Beta) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Beta) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Beta) Reduce(ctx Context) Sort {
	panic("implement me")
}

var _ Sort = Beta{}

const (
	LambdaCmd Name = "λ"
)

func init() {
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) (Context, Sort) {
		err := parseErr(LambdaCmd, []string{"param1", "...", "paramN", "body"}, "where N >= 1")
		if len(list) < 2 {
			panic(err)
		}
		params := make([]Sort, 0, len(list)-1)
		for i := 0; i < len(list)-1; i++ {
			var param Sort

		}

	}

}

type Lambda struct {
	Param Name
	Body  Sort
}

func (l Lambda) Form() Form {
	return List{LambdaCmd, l.Param, l.Body.Form()}
}

func (l Lambda) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (l Lambda) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (l Lambda) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (l Lambda) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (l Lambda) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Lambda{}
