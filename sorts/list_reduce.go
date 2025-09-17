package sorts

const (
	BetaCmd Name = "β"
)

func init() {
	ListParseFuncMap[BetaCmd] = func(ctx Context, list List) Sort {
		err := parseErr(BetaCmd, []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}

		cmd := ctx.Parse(list[0])
		args := make([]Sort, 0, len(list)-1)
		for i := 1; i < len(list); i++ {
			args = append(args, ctx.Parse(list[i]))
		}
		if len(args) == 0 {
			return cmd
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

		return output
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
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := parseErr(LambdaCmd, []string{
			makeForm(AnnotCmd, "param1", "type1"),
			"...",
			makeForm(AnnotCmd, "paramN", "typeN"),
			"body",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}
		params := make([]Annot, 0, len(list)-1)
		for i := 0; i < len(list)-1; i++ {
			params = append(params, mustType[Annot](err, ctx.Parse(list[i])))
		}
		body := ctx.Parse(list[len(list)-1])
		if len(params) == 0 {
			return body
		}

		output := Lambda{
			Param: params[len(params)-1],
			Body:  body,
		}
		for i := len(params) - 2; i >= 0; i-- {
			output = Lambda{
				Param: params[i],
				Body:  output,
			}
		}
		return output
	}

}

type Lambda struct {
	Param Annot
	Body  Sort
}

func (l Lambda) Form() Form {
	return List{LambdaCmd, l.Param.Form(), l.Body.Form()}
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

const (
	InhabitedCmd Name = "inh"
)

func init() {
	ListParseFuncMap[InhabitedCmd] = func(ctx Context, list List) Sort {
		err := parseErr(InhabitedCmd, []string{"type"})
		if len(list) != 1 {
			panic(err)
		}
		t := ctx.Parse(list[0])
		return Inhabited{
			Atom: NewTerm(List{InhabitedCmd, t.Form()}, t),
			Type: t,
		}
	}
}

type Inhabited struct {
	Atom Atom
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{InhabitedCmd, s.Type.Form()}
}

func (s Inhabited) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Inhabited) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Inhabited) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Inhabited) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Inhabited) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Inhabited{}
