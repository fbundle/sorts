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
		return ctx, Beta{
			Cmd:  cmd,
			Args: args,
		}
	}
}

type Beta struct {
	Cmd  Sort
	Args []Sort
}

func (b Beta) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Beta{}

const (
	Lambda Name = "λ"
)
