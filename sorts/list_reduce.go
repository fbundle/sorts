package sorts

import "slices"

func init() {
	DefaultCompileFunc = func(ctx Context, list List) Sort {
		err := compileErr("", []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}

		cmd := ctx.Compile(list[0])
		args := make([]Sort, 0, len(list)-1)
		for i := 1; i < len(list); i++ {
			args = append(args, ctx.Compile(list[i]))
		}
		if len(args) == 0 {
			return cmd
		}

		output := (Beta{
			Cmd: cmd,
			Arg: args[0],
		}).TypeCheck(ctx)
		for i := 1; i < len(args); i++ {
			output = (Beta{
				Cmd: output,
				Arg: args[i],
			}).TypeCheck(ctx)
		}

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

func (s Beta) TypeCheck(ctx Context) Sort {
	s = Beta{
		Cmd: s.Cmd.TypeCheck(ctx),
		Arg: s.Arg.TypeCheck(ctx),
	}

	arrow := mustType[Arrow](TypeErr, s.Cmd.Parent(ctx))
	A := s.Arg.Parent(ctx)

	if !A.LessEqual(ctx, arrow.A) {
		panic(TypeErr)
	}
	return s
}

func (s Beta) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Beta) Parent(ctx Context) Sort {
	arrow := mustType[Arrow](TypeErr, s.Cmd.Parent(ctx))
	return arrow.B
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
	LambdaCmd Name = "=>"
)

func init() {
	ListCompileFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := compileErr(LambdaCmd, []string{
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
			params = append(params, parseAnnot(ctx, list[i]))
		}
		body := ctx.Compile(list[len(list)-1])
		if len(params) == 0 {
			return body
		}

		output := (Lambda{
			Param: params[len(params)-1],
			Body:  body,
		}).TypeCheck(ctx)
		for i := len(params) - 2; i >= 0; i-- {
			output = (Lambda{
				Param: params[i],
				Body:  output,
			}).TypeCheck(ctx)
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

func (l Lambda) TypeCheck(ctx Context) Sort {
	return Lambda{
		Param: l.Param,
		Body:  l.Body.TypeCheck(ctx),
	}
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
	ListCompileFuncMap[InhabitedCmd] = func(ctx Context, list List) Sort {
		err := compileErr(InhabitedCmd, []string{"type"})
		if len(list) != 1 {
			panic(err)
		}
		t := ctx.Compile(list[0])
		return Inhabited{
			Sort: NewTerm(List{InhabitedCmd, t.Form()}, t),
			Type: t,
		}
	}
}

type Inhabited struct {
	Sort Sort
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{InhabitedCmd, s.Type.Form()}
}

func (s Inhabited) TypeCheck(ctx Context) Sort {
	return Inhabited{
		Sort: s.Sort.TypeCheck(ctx),
		Type: s.Type.TypeCheck(ctx),
	}
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

const (
	InductiveCmd = "inductive"
)

func init() {
	ListCompileFuncMap[InductiveCmd] = func(ctx Context, list List) Sort {
		err := compileErr(InductiveCmd, []string{
			"name",
			makeForm(AnnotCmd, "constructor1", "type1"),
			"...",
			makeForm(AnnotCmd, "constructorN", "typeN"),
		}, "where N >= 1")
		if len(list) < 2 {
			panic(err)
		}
		name := mustType[Name](err, list[0])
		subCtx := ctx.Set(name, nil)
		mks := make([]Annot, 0, len(list)-1)
		for i := 1; i < len(list); i++ {
			mks = append(mks, parseAnnot(subCtx, list[i]))
		}

		return Inductive{
			Name: name,
			Mks:  mks,
		}
	}
}

type Inductive struct {
	Name Name
	Mks  []Annot
}

func (s Inductive) Form() Form {
	form := List{InhabitedCmd, s.Name}
	for _, mk := range s.Mks {
		form = append(form, mk.Form())
	}
	return form
}

func (s Inductive) TypeCheck(ctx Context) Sort {
	s = Inductive{
		Name: s.Name,
		Mks:  ,
	}


}

func (s Inductive) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Inductive) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Inductive) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Inductive) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Inductive{}

type Match struct {
	Cond  Inductive
	Cases []Case
}
