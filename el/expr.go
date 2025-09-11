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

func ParseForm(e form.Form) (Expr, error) {
	switch e := e.(type) {
	case form.Term:
		return Term(e), nil
	case form.List:
		if len(e) == 0 {
			return nil, errors.New("empty list")
		}
		head, list := e[0], e[1:]

		// Is it a special form?
		if cmdTerm, ok := head.(form.Term); ok {
			if parser, ok := globalListParsers[cmdTerm]; ok {
				return parser(ParseForm, list)
			}
		}

		// It's a regular function call
		cmd, err := ParseForm(head)
		if err != nil {
			return nil, err
		}
		args := make([]Expr, len(list))
		for i, arg := range list {
			args[i], err = ParseForm(arg)
			if err != nil {
				return nil, err
			}
		}
		return FunctionCall{Cmd: cmd, Args: args}, nil
	default:
		return nil, errors.New("unknown form")
	}
}
