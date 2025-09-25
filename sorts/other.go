package sorts

const (
	LetCmd Name = "let"
)

type Let struct {
	Binding Binding
	Body    Sort
}

func (s Let) Form() Form {
	return List{LetCmd, s.Binding.Form(), s.Body.Form()}
}

func (s Let) Parent(ctx Context) Sort {
	return s.Body.Parent(ctx.Set(s.Binding.Name, s.Binding.Value))
}

func (s Let) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Let) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Let) Eval(ctx Context) Sort {
	return s.Body.Eval(ctx.Set(s.Binding.Name, s.Binding.Value))
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
