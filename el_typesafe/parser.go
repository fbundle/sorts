package el_typesafe

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type ParseFunc = func(form.Form) (Expr, error)
type ListParseFunc = func(ParseFunc, form.List) (Expr, error)

var defaultParser parser

func ParseForm(e form.Form) (Expr, error) {
	return defaultParser.parseForm(e)
}
func RegisterListParser(cmd form.Name, listParser func(ParseFunc, form.List) (Expr, error)) {
	defaultParser = defaultParser.registerListParser(cmd, listParser)
}

type parser struct {
	listParsers map[form.Name]ListParseFunc
}

func (parser parser) registerListParser(cmd form.Name, listParser func(ParseFunc, form.List) (Expr, error)) parser {
	if parser.listParsers == nil {
		parser.listParsers = make(map[form.Name]ListParseFunc)
	}
	parser.listParsers[cmd] = listParser
	return parser
}

func (parser parser) parseForm(e form.Form) (Expr, error) {
	switch e := e.(type) {
	case form.Name:
		return Term(e), nil
	case form.List:
		if len(e) == 0 {
			return nil, errors.New("empty list")
		}
		head, list := e[0], e[1:]

		// Is it a special form?
		if cmdTerm, ok := head.(form.Name); ok {
			if listParser, ok := parser.listParsers[cmdTerm]; ok {
				return listParser(parser.parseForm, list)
			}
		}

		// It's a regular function call
		// convert (f a b c) into (((f a) b) c)

		cmd, err := parser.parseForm(head)
		if err != nil {
			return nil, err
		}
		for _, argExpr := range list {
			arg, err := parser.parseForm(argExpr)
			if err != nil {
				return nil, err
			}
			cmd = FunctionCall{Cmd: cmd, Arg: arg}
		}
		return cmd, nil
	default:
		return nil, errors.New("unknown form")
	}
}
