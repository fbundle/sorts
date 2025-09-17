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

func (b Beta) Reduce(ctx Context) Sort {
	panic("implement me")
}

var _ Sort = Beta{}

const (
	LambdaCmd Name = "λ"
)

type Lambda struct {
	Params []Name
	Body   Sort
}

func (l Lambda) Form() Form {
	form := List{}
	form = append(form, LambdaCmd)
	for _, param := range l.Params {
		form = append(form, param)
	}
	form = append(form, l.Body.Form())
	return form
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
