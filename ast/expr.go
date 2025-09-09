package ast

type Token = string

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
	output = append(output, TokenBlockBegin)
	for _, arg := range e {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, TokenBlockEnd)
	return output
}
