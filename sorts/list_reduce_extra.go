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

func (s Annot) inhabited() Inhabited {
	return Inhabited{
		Name: s.Name,
		Type: s.Type,
	}
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

func (c Case) Form() Form {
	if len(c.MkArgs) == 0 {
		return List{Name(CaseCmd), c.MkName, c.Value.Form()}
	} else {
		return List{Name(CaseCmd), append(List{c.MkName}, slicesMap(c.MkArgs, func(name Name) Form {
			return name
		})...), c.Value.Form()}
	}
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

	var c Case
	switch cList := list[0].(type) {
	case Name:
		c = Case{
			MkName: cList,
			MkArgs: nil,
			Value:  ctx.Compile(list[1]).TypeCheck(ctx),
		}
	case List:
		c = Case{
			MkName: mustType[Name](err, cList[0]),
			MkArgs: slicesMap(cList[1:], func(form Form) Name {
				return mustType[Name](err, form)
			}),
			Value: ctx.Compile(list[1]).TypeCheck(ctx),
		}
	default:
		panic(err)
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
