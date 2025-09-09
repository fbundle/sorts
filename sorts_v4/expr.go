package sorts

type Token string

const (
	TokenOpen  Token = "("
	TokenClose Token = ")"
)

type Expr interface {
	Marshal() []Token
}

type Term string

func (e Term) Marshal() []Token {
	return []Token{Token(e)}
}

type Node struct {
	Cmd  Term
	Args []Expr
}

func (e Node) Marshal() []Token {
	var output []Token
	output = append(output, TokenOpen)
	output = append(output, Token(e.Cmd))
	for _, arg := range e.Args {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, TokenClose)
	return output
}
