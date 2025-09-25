package sorts

// TODO finish later

const (
	LetCmd Name = "let"
)

type Let struct {
	Binding Binding
	Body    Sort
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
