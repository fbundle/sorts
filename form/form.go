package form

import "strings"

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

const (
	BlockBeg Token = "("
	BlockEnd Token = ")"
)

func (n List) Marshal() []Token {
	var output []Token
	output = append(output, BlockBeg)
	for _, arg := range n {
		output = append(output, arg.Marshal()...)
	}
	output = append(output, BlockEnd)
	return output
}

func (n List) formAttr() {}

func String(form Form) string {
	s := strings.Join(form.Marshal(), " ")
	s = strings.ReplaceAll(s, "( ", "(")
	s = strings.ReplaceAll(s, " )", ")")
	return s
}
