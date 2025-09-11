package el

import "github.com/fbundle/sorts/form"

func ParseForm(e form.Form) (Expr, error) {
	return defaultParser.Parse(e)
}

var defaultParser = Parser{
	ListParsers: map[form.Term]ListParseFunc{
		"=>":    parseLambda,
		":":     parseDefine,
		":=":    parseAssign,
		"chain": parseChain,
		"match": parseMatch,
	},
}

func parseLambda(parse ParseFunc, list form.List) (Expr, error) {
	panic("implement me")
}
func parseDefine(parse ParseFunc, list form.List) (Expr, error) {
	panic("implement me")
}
func parseAssign(parse ParseFunc, list form.List) (Expr, error) {
	panic("implement me")
}
func parseChain(parse ParseFunc, list form.List) (Expr, error) {
	panic("implement me")
}
func parseMatch(parse ParseFunc, list form.List) (Expr, error) {
	panic("implement me")
}
