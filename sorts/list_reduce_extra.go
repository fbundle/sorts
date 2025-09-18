package sorts

const (
	AnnotCmd Name = ":"
)

func compileAnnot(ctx Context, form Form) Annot {
	err := compileErr(AnnotCmd, []string{"name", "type"})
	list := mustType[List](err, form)
	if len(list) != 2 {
		panic(err)
	}

	return Annot{
		Name: mustType[Name](err, list[0]),
		Type: ctx.Compile(list[1]).TypeCheck(ctx),
	}
}

type Annot struct {
	Name Name
	Type Sort
}

func (s Annot) Form() Form {
	return List{AnnotCmd, s.Name, s.Type.Form()}
}

const (
	CaseCmd   Name = "=>"
	CaseFinal Name = "_"
)

type Case struct {
	MkName Name
	MkArgs []Name
	Value  Sort
}

func compileCase(ctx Context, form Form) Case {
	err := compileErr(CaseCmd, []string{
		makeForm("constructor", "arg1", "...", "argN"),
		"value",
	})
	list := mustType[List](err, form)
	if len(list) != 2 {
		panic(err)
	}

	cList := mustType[List](err, list[0])
	return Case{
		MkName: mustType[Name](err, cList[0]),
		MkArgs: slicesMap(cList[1:], func(form Form) Name {
			return mustType[Name](err, form)
		}),
		Value: ctx.Compile(list[1]).TypeCheck(ctx),
	}
}
