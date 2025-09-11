package form

type Token = string

type Form interface {
	Marshal() []Token
	formAttr()
}

type Term string

func (t Term) Marshal() []Token {
	return []Token{Token(t)}
}

func (t Term) formAttr() {}

type List []Form

func (n List) Marshal() []Token {
	var output []Token
	output = append(output, TokenBlockBegin)
	for _, arg := range n {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, TokenBlockEnd)
	return output
}

func (n List) formAttr() {}
