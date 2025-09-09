package lisp_util

type Token = string

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

type Node []Expr

func (e Node) Marshal() []Token {
	var output []Token
	output = append(output, TokenOpen)
	for _, arg := range e {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, TokenClose)
	return output
}
