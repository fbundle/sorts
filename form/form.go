package form

import "strings"

type Token = string

// Form - Union[Name, List]
type Form interface {
	Marshal() []Token
	formAttr()
}

type Name string

func (t Name) Marshal() []Token {
	return []Token{Token(t)}
}

func (t Name) formAttr() {}

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
