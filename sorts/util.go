package sorts

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

func (a Annot) Push(ctx Context) Context {
	return ctx.Set(a.Name, NewTerm(a.Name, a.Type))
}
