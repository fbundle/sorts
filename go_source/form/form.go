package form

import "strings"

type Token = string

const (
	blockBeg Token = "("
	blockEnd Token = ")"
)

// Form - Union[Cmd, List]
type Form interface {
	Marshal() []Token
	formAttr()
}
type Name string
type List []Form

func (t Name) Marshal() []Token {
	return []Token{Token(t)}
}

func (t Name) formAttr() {}

func (n List) Marshal() []Token {
	var output []Token
	output = append(output, blockBeg)
	for _, arg := range n {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, blockEnd)
	return output
}

func (n List) formAttr() {}

func String(form Form) string {
	return strings.Join(form.Marshal(), " ")
}
