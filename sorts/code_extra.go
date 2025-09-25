package sorts

const (
	LetCmd Name = "let"
)

type Let struct {
	Binding Binding
	Body    Code
}

func (c Let) Form() Form {
	return List{LetCmd, c.Binding.Form(), c.Body.Form()}
}

func (c Let) Eval(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
