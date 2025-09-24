package sorts

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

const (
	AnnotCmd Name = ":"
)

type Annot struct {
	Name Name
	Type Sort
}

func (a Annot) Form() Form {
	return List{AnnotCmd, a.Name, a.Type.Form()}
}
