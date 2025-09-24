package sorts

const (
	TypeCmd Name = "&"
)

type Type struct {
	Sort
	Body Sort
}

func (s Type) Form() Form {
	return List{TypeCmd, s.Body.Form()}
}

func (s Type) Compile(ctx Context) Sort {
	s.Sort = s.Body.Parent(ctx)
	return s
}
