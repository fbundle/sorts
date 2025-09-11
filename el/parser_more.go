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

func parseLambda(list []form.Form) (Expr, error) {
	panic("implement me")
}
func parseDefine(list []form.Form) (Expr, error) {
	panic("implement me")
}
func parseAssign(list []form.Form) (Expr, error) {
	panic("implement me")
}
func parseMatch(list []form.Form) (Expr, error) {
	panic("implement me")
}
func parseChain(list []form.Form) (Expr, error) {
	panic("implement me")
}
