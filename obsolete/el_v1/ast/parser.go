package ast

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

func Parse(e form.Form) (Expr, error) {
	switch e := e.(type) {
	case form.Term:
		return Name(e), nil
	case form.List:
		head, args := e[0], e[1:]
		cmd, ok := head.(form.Term)
		if !ok {
			return nil, errors.New("expected term expression")
		}
		switch cmd {
		case "let":
			if len(args)%2 == 0 {
				return nil, errors.New("expected at least even number of args")
			}
			var bindings []LetBinding
			for i := 0; i <= len(args)-3; i += 2 {
				name, ok := args[i].(form.Term)
				if !ok {
					return nil, errors.New("expected term expression")
				}
				value, err := Parse(args[i+1])
				if err != nil {
					return nil, err
				}
				bindings = append(bindings, LetBinding{
					Name: Name(name),
					Expr: value,
				})
			}
			final, err := Parse(args[len(args)-1])
			if err != nil {
				return nil, err
			}
			return Let{
				Bindings: bindings,
				Final:    final,
			}, nil
		case "match":
			if len(args)%2 == 1 {
				return nil, errors.New("expected at least even number of args")
			}
			cond, err := Parse(args[0])
			if err != nil {
				return nil, err
			}
			var cases []MatchCase
			for i := 1; i <= len(args)-3; i += 2 {
				comp, err := Parse(args[i])
				if err != nil {
					return nil, err
				}
				value, err := Parse(args[i+1])
				if err != nil {
					return nil, err
				}
				cases = append(cases, MatchCase{
					Comp:  comp,
					Value: value,
				})
			}
			final, err := Parse(args[len(args)-1])
			if err != nil {
				return nil, err
			}
			return Match{
				Cond:    cond,
				Cases:   cases,
				Default: final,
			}, nil
		case form.TermArrowDouble: // lambda
			if len(args) != 2 {
				return nil, errors.New("expected 2 args")
			}
			param, ok := args[0].(form.Term)
			if !ok {
				return nil, errors.New("expected term expression")
			}
			body, err := Parse(args[1])
			if err != nil {
				return nil, err
			}
			return Lambda{
				Param: Name(param),
				Body:  body,
			}, nil
		default:
			var aList []Expr
			for _, arg := range args {
				a, err := Parse(arg)
				if err != nil {
					return nil, err
				}
				aList = append(aList, a)
			}
			return FunctionCall{
				Cmd:  Name(cmd),
				Args: aList,
			}, nil
		}
	default:
		return nil, errors.New("invalid expression")
	}
}
