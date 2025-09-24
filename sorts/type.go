package sorts

const (
	TypeCmd Name = "&"
)

type Type struct {
	Body Sort
}

func (s Type) Form() Form {
	return List{TypeCmd, s.Body.Form()}
}

func (s Type) Parent(ctx Context) Sort {
	return s.Body.Parent(ctx).Parent(ctx)
}
