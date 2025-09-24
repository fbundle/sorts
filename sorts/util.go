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
