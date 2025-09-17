package sorts

const (
	AnnotCmd Name = ":"
)

func init() {
	ListParseFuncMap[AnnotCmd] = func(ctx Context, list List) Sort {

	}
}

type Annot struct {
	Name Name
	Type Sort
}

func (a Annot) Form() Form {
	return List{AnnotCmd, a.Name, a.Type.Form()}
}

var _ Code = Annot{}
