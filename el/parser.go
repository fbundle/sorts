package el

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type ParseFunc = func(form.Form) (Expr, error)
type ListParseFunc = func(ParseFunc, form.List) (Expr, error)

type Parser struct {
	ListParsers map[form.Term]ListParseFunc
}

func parseFunctionCall(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) == 0 {
		return nil, errors.New("empty list")
	}
	cmd, err := parse(list[0])
	if err != nil {
		return nil, err
	}
	args := make([]Expr, 0, len(list)-1)
	for i := 1; i < len(list); i++ {
		arg, err := parse(list[i])
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return FunctionCall{
		Cmd:  cmd,
		Args: args,
	}, nil
}

func (parser Parser) Parse(e form.Form) (Expr, error) {
	switch e := e.(type) {
	case form.Term:
		return Term(e), nil
	case form.List:
		head, args := e[0], e[1:]
		listParseFunc := func() ListParseFunc {
			cmd, ok := head.(form.Term)
			if !ok {
				return parseFunctionCall
			}
			listParseFunc, ok := parser.ListParsers[cmd]
			if !ok {
				return parseFunctionCall
			}
			return listParseFunc
		}()
		return listParseFunc(parser.Parse, args)
	default:
		return nil, errors.New("unknown form")
	}
}
