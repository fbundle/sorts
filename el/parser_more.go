package el

import "github.com/fbundle/sorts/form"

func Parse(e form.Form) (Expr, error) {
	return defaultParser.Parse(e)
}

var defaultParser = Parser{
	ListParsers: map[form.Term]ListParseFunc{
		"=>":    nil,
		":":     nil,
		":=":    nil,
		"chain": nil,
		"match": nil,
	},
}
