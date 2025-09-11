package el

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

type Expr interface {
	Marshal() form.Form
	mustExpr()
}

type Term string

func (t Term) mustExpr() {}

func (t Term) Marshal() form.Form {
	return form.Term(t)
}

// FunctionCall - (cmd arg1 arg2 ...)
type FunctionCall struct {
	Cmd  Expr
	Args []Expr
}

func (f FunctionCall) mustExpr() {}

func (f FunctionCall) Marshal() form.Form {
	forms := make([]form.Form, 0, 1+len(f.Args))
	forms = append(forms, f.Cmd.Marshal())
	for _, arg := range f.Args {
		forms = append(forms, arg.Marshal())
	}
	return form.List(forms)
}

type ParseFunc = func(form.Form) (Expr, error)
type ListParseFunc = func(ParseFunc, form.List) (Expr, error)

var globalListParsers = make(map[form.Term]ListParseFunc)

func RegisterListParser(cmd form.Term, listParser func(ParseFunc, form.List) (Expr, error)) {
	globalListParsers[cmd] = listParser
}

func parseFunctionCall(parseFunc ParseFunc, list form.List) (Expr, error) {
	if len(list) == 0 {
		return nil, errors.New("empty list")
	}
	cmd, err := parseFunc(list[0])
	if err != nil {
		return nil, err
	}
	args := make([]Expr, 0, len(list)-1)
	for i := 1; i < len(list); i++ {
		arg, err := parseFunc(list[i])
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

func ParseForm(e form.Form) (Expr, error) {
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
			listParseFunc, ok := globalListParsers[cmd]
			if !ok {
				return parseFunctionCall
			}
			return listParseFunc
		}()
		return listParseFunc(ParseForm, args)
	default:
		return nil, errors.New("unknown form")
	}
}
