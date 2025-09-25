package sorts

const (
	AnnotCmd Name = ":"
)

func compileAnnot(ctx Context, list List) Annot {
	err := compileErr(list, []string{string(AnnotCmd), "name", "type"})
	if len(list) != 2 {
		panic(err)
	}
	return Annot{
		Name: mustType[Name](err, list[0]),
		Type: ctx.Parse(list[1]),
	}
}

type Annot struct {
	Name Name
	Type Code
}

func (a Annot) Form() Form {
	return List{AnnotCmd, a.Name, a.Type.Form()}
}

const (
	BindingCmd Name = ":="
)

func compileBinding(ctx Context, list List) Binding {
	err := compileErr(list, []string{string(BindingCmd), "name", "value"})
	if len(list) != 2 {
		panic(err)
	}
	return Binding{
		Name:  mustType[Name](err, list[0]),
		Value: ctx.Parse(list[1]),
	}
}

type Binding struct {
	Name  Name
	Value Code
}

func (b Binding) Form() Form {
	return List{BindingCmd, b.Name, b.Value.Form()}
}
