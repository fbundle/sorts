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
		"chain": parseMatch,
		"match": parseChain,
	},
}

func parseLambda(list form.List) (Expr, error) {
	panic("implement me")
}
func parseDefine(list form.List) (Expr, error) {
	panic("implement me")
}
func parseAssign(list form.List) (Expr, error) {
	panic("implement me")
}
func parseMatch(list form.List) (Expr, error) {
	panic("implement me")
}
func parseChain(list form.List) (Expr, error) {
	panic("implement me")
}
