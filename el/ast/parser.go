package ast

import (
	"errors"

	"github.com/fbundle/sorts/expr"
)

func Parse(e expr.Expr) (Expr, error) {
	switch e := e.(type) {
	case expr.Term:
		return Name(e), nil
	case expr.Node:
		head, args := e[0], e[1:]
		cmd, ok := head.(expr.Term)
		if !ok {
			return nil, errors.New("expected term expression")
		}
		switch cmd {
		case "let":
			if len(args)%2 == 0 {
				return nil, errors.New("expected at least even number of args")
			}
			var bindings []Binding
			for i := 0; i <= len(args)-3; i += 2 {
				name, ok := args[i].(expr.Term)
				if !ok {
					return nil, errors.New("expected term expression")
				}
				value, err := Parse(args[i+1])
				if err != nil {
					return nil, err
				}
				bindings = append(bindings, Binding{
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
			cond := args[0]
			var cases []Case
			for i := 1; i <= len(args)-3; i += 2 {
				comp, err := Parse(args[i])
				if err != nil {
					return nil, err
				}
				value, err := Parse(args[i+1])
				if err != nil {
					return nil, err
				}
				cases = append(cases, Case{
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
		case expr.TermArrowDouble: // lambda
			if len(args) != 2 {
				return nil, errors.New("expected 2 args")
			}
			param, ok := args[0].(expr.Term)
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
			return Call{
				Cmd:  Name(cmd),
				Args: aList,
			}, nil
		}
	default:
		return nil, errors.New("invalid expression")
	}
}
