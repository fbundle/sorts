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
	return c.Body.Eval(ctx.Set(c.Binding.Name, c.Binding.Value.Eval(ctx)))
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
