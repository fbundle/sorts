package sorts

import "fmt"

func init() {
	DefaultCompileFunc = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{"beta", "arg1", "...", "argN"}, "where N >= 0")
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
		fmt.Println("A", arrow.A.Form())
		fmt.Println("a", A.Form())
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
	LambdaCmd = "=>"
)

func init() {
	ListCompileFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{
			LambdaCmd,
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
	return List{Name(LambdaCmd), l.Param.Form(), l.Body.Form()}
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
	InhabitedCmd = "inh"
)

func init() {
	ListCompileFuncMap[InhabitedCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{InhabitedCmd, "type"})
		if len(list) != 1 {
			panic(err)
		}
		t := ctx.Compile(list[0])
		return Inhabited{
			Sort: NewTerm(Inhabited{
				Type: t,
			}.Form(), t),
			Type: t,
		}
	}
}

type Inhabited struct {
	Sort Sort
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{Name(InhabitedCmd), s.Type.Form()}
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
	return s.Type
}

func (s Inhabited) LessEqual(ctx Context, d Sort) bool {
	return false
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
		err := compileErr(list, []string{
			InductiveCmd,
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

func (s Inductive) constructors() map[Name]Sort {
	c := make(map[Name]Sort)
	for _, mk := range s.Mks {
		c[mk.Name] = mk.Type
	}
	return c
}

func (s Inductive) Form() Form {
	return append(List{Name(InhabitedCmd), s.Name}, slicesMap(s.Mks, func(mk Annot) Form {
		return mk.Form()
	})...)
}

func (s Inductive) TypeCheck(ctx Context) Sort {
	s = Inductive{
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
	if len(s.constructors()) != len(s.Mks) {
		panic(TypeErr)
	}
	return s
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

const (
	MatchCmd = "match"
)

func init() {
	ListCompileFuncMap[MatchCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{
			MatchCmd,
			"cond",
			makeForm(CaseCmd, makeForm("constructor1", "arg11", "...", "arg1N"), "value1"),
			"...",
			makeForm(CaseCmd, makeForm("constructorM", "argM1", "...", "argMN"), "valueM"),
		})

		if len(list) < 1 {
			panic(err)
		}

		return Match{
			Cond: mustType[Inductive](err, ctx.Compile(list[0]).TypeCheck(ctx)),
			Cases: slicesMap(list[1:], func(form Form) Case {
				return compileCase(ctx, mustType[List](err, form)[1:])
			}),
		}
	}
}

type Match struct {
	Cond  Inductive
	Cases []Case
}

func (s Match) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (s Match) TypeCheck(ctx Context) Sort {
	condConstructors := s.Cond.constructors()
	caseConstructors := make(map[Name]Case)
	slicesForEach(s.Cases, func(c Case) {
		caseConstructors[c.MkName] = c
	})
	if len(caseConstructors) != len(s.Cases) {
		panic(TypeErr)
	}
	for cname, c := range caseConstructors {
		constr, ok := condConstructors[cname]
		if !ok {
			panic(TypeErr)
		}

		if len(serialize(constr)) != len(c.MkArgs)+1 {
			panic(TypeErr)
		}
	}
	// must match all cases
	if _, ok := caseConstructors[CaseFinal]; !ok {
		// don't have final case
		if len(caseConstructors) != len(condConstructors) {
			panic(TypeErr)
		}
	}
	return s
}

func (s Match) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Match) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Match) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Match) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Match{}

const (
	LetCmd = "let"
)

func init() {
	ListCompileFuncMap[LetCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{
			LetCmd,
			makeForm(BindingCmd, "name1", "value1"),
			"...",
			makeForm(BindingCmd, "nameN", "valueN"),
			"final",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}

		return Let{
			Bindings: slicesMap(list[:len(list)-1], func(form Form) Binding {
				binding := compileBinding(ctx, mustType[List](err, form)[1:])
				ctx = ctx.Set(binding.Name, binding.Value)
				return binding
			}),
			Final: ctx.Compile(list[len(list)-1]).TypeCheck(ctx),
		}
	}
}

type Let struct {
	Bindings []Binding
	Final    Sort
}

func (l Let) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (l Let) TypeCheck(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (l Let) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (l Let) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (l Let) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (l Let) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Let{}
