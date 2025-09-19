package sorts

const (
	AnnotCmd = ":"
)

func compileAnnot(ctx Context, list List) Annot {
	err := compileErr(list, []string{AnnotCmd, "name", "type"})
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
	return List{Name(AnnotCmd), s.Name, s.Type.Form()}
}

const (
	CaseCmd   = "=>"
	CaseFinal = "_"
)

type Case struct {
	MkName Name
	MkArgs []Name
	Value  Sort
}

func compileCase(ctx Context, list List) Case {
	err := compileErr(list, []string{
		CaseCmd,
		makeForm("constructor", "arg1", "...", "argN"),
		"value",
	})
	if len(list) != 2 {
		panic(err)
	}

	cList := mustType[List](err, list[0])
	c := Case{
		MkName: mustType[Name](err, cList[0]),
		MkArgs: slicesMap(cList[1:], func(form Form) Name {
			return mustType[Name](err, form)
		}),
		Value: ctx.Compile(list[1]).TypeCheck(ctx),
	}

	if c.MkName == CaseFinal && len(c.MkArgs) > 0 {
		panic(TypeErr)
	}

	return c
}

const (
	BindingCmd = ":="
)

type Binding struct {
	Name  Name
	Value Sort
}

func compileBinding(ctx Context, list List) Binding {
	err := compileErr(list, []string{BindingCmd, "name", "value"})
	if len(list) != 2 {
		panic(err)
	}
	return Binding{
		Name:  mustType[Name](err, list[0]),
		Value: ctx.Compile(list[1]).TypeCheck(ctx),
	}
}
