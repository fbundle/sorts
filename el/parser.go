package el

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type ListParseFunc = func(form.List) (Expr, error)

type Parser struct {
	ListParsers map[form.Term]ListParseFunc
}

func (parser Parser) defaultListParseFunc(list form.List) (Expr, error) {
	if len(list) == 0 {
		return nil, errors.New("empty list")
	}
	cmd, err := parser.Parse(list[0])
	if err != nil {
		return nil, err
	}
	args := make([]Expr, 0, len(list)-1)
	for i := 1; i < len(list); i++ {
		arg, err := parser.Parse(list[i])
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
				return parser.defaultListParseFunc
			}
			listParseFunc, ok := parser.ListParsers[cmd]
			if !ok {
				return parser.defaultListParseFunc
			}
			return listParseFunc
		}()
		return listParseFunc(args)
	default:
		return nil, errors.New("unknown form")
	}
}
