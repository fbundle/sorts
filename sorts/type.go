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

func (s Type) Level(ctx Context) int {
	return s.Body.Level(ctx) + 1
}
func (s Type) LessEqual(ctx Context, d Sort) bool {
	return s.Body.Parent(ctx).LessEqual(ctx, d)
}

func (s Type) Reduce(ctx Context) Sort {
	return s.Body.Parent(ctx)
}
