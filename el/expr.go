package el

import (
	"errors"
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func String(e Expr) string {
	return form.String(e.Marshal())
}

type Expr interface {
	Marshal() form.Form
	Resolve(frame Frame) (Frame, sorts.Sort, Expr, error)
	mustExpr()
}

type Term string

func (t Term) mustExpr() {}

func (t Term) Marshal() form.Form {
	return form.Term(t)
}

func (t Term) Resolve(frame Frame) (Frame, sorts.Sort, Expr, error) {
	sort, next, err := frame.get(t)
	if err != nil {
		return frame, nil, nil, err
	}
	if next == t {
		// term not assigned, set as cycle
		return frame, sort, next, nil
	}
	return next.Resolve(frame)
}

// FunctionCall - (cmd arg1 arg2 ...)
type FunctionCall struct {
	Cmd Expr
	Arg Expr
}

func (f FunctionCall) mustExpr() {}

func (f FunctionCall) Marshal() form.Form {
	return form.List{
		f.Cmd.Marshal(),
		f.Arg.Marshal(),
	}
}
func (f FunctionCall) Resolve(frame Frame) (Frame, sorts.Sort, Expr, error) {
	frame, argSort, argValue, err := f.Arg.Resolve(frame)
	if err != nil {
		return frame, nil, nil, err
	}
	frame, cmdSort, cmdValue, err := f.Cmd.Resolve(frame)
	if err != nil {
		return frame, nil, nil, err
	}
	switch cmd := cmdValue.(type) {
	case Lambda:
		frame, err := frame.set(cmd.Param, argSort, argValue)
		if err != nil {
			return frame, nil, nil, err
		}
		return cmd.Body.Resolve(frame)
	default:
		if B, ok := frame.typeCheckFunctionCall(cmdSort, argSort); ok {
			return frame, sorts.NewTerm(B, fmt.Sprintf("(%s %s)", String(cmd), String(argValue))), FunctionCall{cmd, argValue}, nil
		}
		return frame, nil, nil, fmt.Errorf("type_error: cmd %s, arg %s", sorts.Name(cmdSort), sorts.Name(argSort))
	}
}

type ParseFunc = func(form.Form) (Expr, error)
type ListParseFunc = func(ParseFunc, form.List) (Expr, error)

var defaultParser parser

func ParseForm(e form.Form) (Expr, error) {
	return defaultParser.parseForm(e)
}
func RegisterListParser(cmd form.Term, listParser func(ParseFunc, form.List) (Expr, error)) {
	defaultParser = defaultParser.registerListParser(cmd, listParser)
}

type parser struct {
	listParsers map[form.Term]ListParseFunc
}

func (parser parser) registerListParser(cmd form.Term, listParser func(ParseFunc, form.List) (Expr, error)) parser {
	if parser.listParsers == nil {
		parser.listParsers = make(map[form.Term]ListParseFunc)
	}
	parser.listParsers[cmd] = listParser
	return parser
}

func (parser parser) parseForm(e form.Form) (Expr, error) {
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
