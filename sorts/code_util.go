package sorts

const (
	AnnotCmd Name = ":"
)

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

type Binding struct {
	Name  Name
	Value Code
}

func (b Binding) Form() Form {
	return List{BindingCmd, b.Name, b.Value.Form()}
}
