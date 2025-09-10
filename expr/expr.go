package expr

type Token = string

type Expr interface {
	Marshal() []Token
	exprAttr()
}

type Term string

func (t Term) Marshal() []Token {
	return []Token{Token(t)}
}

func (t Term) exprAttr() {}

type Node []Expr

func (n Node) Marshal() []Token {
	var output []Token
	output = append(output, TokenBlockBegin)
	for _, arg := range n {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, TokenBlockEnd)
	return output
}

func (n Node) exprAttr() {}
