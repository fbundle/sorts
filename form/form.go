package form

type Token = string

// Form - Union[Name, List]
type Form interface {
	Marshal(blockBeg Token, blockEnd Token) []Token
	formAttr()
}
type Name string
type List []Form

func (t Name) Marshal(blockBeg Token, blockEnd Token) []Token {
	return []Token{Token(t)}
}

func (t Name) formAttr() {}

func (n List) Marshal(blockBeg Token, blockEnd Token) []Token {
	var output []Token
	output = append(output, blockBeg)
	for _, arg := range n {
		output = append(output, arg.Marshal(blockBeg, blockEnd)...)
	}
	output = append(output, blockEnd)
	return output
}

func (n List) formAttr() {}
