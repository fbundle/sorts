package sorts

const (
	AnnotCmd Name = ":"
)

func parseAnnot(ctx Context, form Form) Annot {
	err := compileErr(AnnotCmd, []string{"name", "type"})
	list := mustType[List](err, form)
	if len(list) != 2 {
		panic(err)
	}

	return Annot{
		Name: mustType[Name](err, list[0]),
		Type: ctx.Compile(list[1]),
	}
}

type Annot struct {
	Name Name
	Type Sort
}

func (s Annot) Form() Form {
	return List{AnnotCmd, s.Name, s.Type.Form()}
}

var _ Code = Annot{}

const (
	CaseCmd Name = "case"
)

type Case struct {
	MkName Name
	MkArgs []Name
	Value  Sort
}
