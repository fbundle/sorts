package sorts

func init() {
	DefaultCompileFunc = func(ctx Context, list List) Sort {
		err := compileErr("", []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}

		cmd := ctx.Compile(list[0]).TypeCheck(ctx)
		args := slicesMap(list[1:], func(form Form) Sort {
			return ctx.Compile(form).TypeCheck(ctx)
		})

		return slicesReduce(args, cmd, func(output Sort, arg Sort) Sort {
			return (Beta{
				Cmd: output,
				Arg: arg,
			}).TypeCheck(ctx)
		})
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

		params := slicesMap(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(ctx, form)
		})
		body := ctx.Compile(list[len(list)-1]).TypeCheck(ctx)

		return slicesReduce(slicesReverse(params), body, func(output Sort, param Annot) Sort {
			return (Lambda{
				Param: param,
				Body:  output,
			}).TypeCheck(ctx)
		})
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

		mks := slicesMap(list[1:], func(form Form) Annot {
			return compileAnnot(subCtx, form)
		})

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
	return append(List{InhabitedCmd, s.Name}, slicesMap(s.Mks, func(mk Annot) Form {
		return mk.Form()
	})...)
}

func (s Inductive) TypeCheck(ctx Context) Sort {
	return Inductive{
		Name: s.Name,
		Mks: slicesMap(s.Mks, func(mk Annot) Annot {
			t := mk.Type.TypeCheck(ctx)
			arrow := serialize(t)
			if arrow[len(arrow)-1].Form() != s.Name {
				panic(TypeErr)
			}
			return Annot{
				Name: mk.Name,
				Type: t,
			}
		}),
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

func (m Match) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (m Match) TypeCheck(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (m Match) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (m Match) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (m Match) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (m Match) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Match{}

type Let struct{}
