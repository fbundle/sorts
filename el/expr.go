package el

import (
	"errors"
	"fmt"

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

var defaultParser Parser

func ParseForm(e form.Form) (Expr, error) {
	return defaultParser.ParseForm(e)
}
func RegisterListParser(cmd form.Term, listParser func(ParseFunc, form.List) (Expr, error)) {
	defaultParser.RegisterListParser(cmd, listParser)
}

type Parser struct {
	listParsers map[form.Term]ListParseFunc
}

func (parser Parser) RegisterListParser(cmd form.Term, listParser func(ParseFunc, form.List) (Expr, error)) {
	if parser.listParsers == nil {
		parser.listParsers = make(map[form.Term]ListParseFunc)
	}
	parser.listParsers[cmd] = listParser
	fmt.Printf("Registered list parser for %s\n", cmd)
}

func (parser Parser) ParseForm(e form.Form) (Expr, error) {
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
			if listParser, ok := parser.listParsers[cmdTerm]; ok {
				return listParser(parser.ParseForm, list)
			}
		}

		// It's a regular function call
		cmd, err := parser.ParseForm(head)
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
