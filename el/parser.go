package el

import (
	"errors"
	"fmt"

	"github.com/fbundle/sorts/form"
)

// ParseForm converts a form.Form into an el.Expr
func ParseForm(f form.Form) (Expr, error) {
	switch v := f.(type) {
	case form.Term:
		return Term(v), nil
	case form.List:
		return parseList(v)
	default:
		return nil, fmt.Errorf("unknown form type: %T", f)
	}
}

func parseList(l form.List) (Expr, error) {
	if len(l) == 0 {
		return nil, errors.New("empty list")
	}

	// First element should be a term (the command/operator)
	cmd, ok := l[0].(form.Term)
	if !ok {
		return nil, errors.New("first element of list must be a term")
	}

	args := l[1:]

	switch cmd {
	case "=>":
		if len(args) != 2 {
			return nil, errors.New("lambda must have exactly 2 arguments: param and body")
		}
		param, ok := args[0].(form.Term)
		if !ok {
			return nil, errors.New("lambda parameter must be a term")
		}
		body, err := ParseForm(args[1])
		if err != nil {
			return nil, fmt.Errorf("parsing lambda body: %w", err)
		}
		return Lambda{
			Param: Term(param),
			Body:  body,
		}, nil

	case ":":
		if len(args) != 2 {
			return nil, errors.New("define must have exactly 2 arguments: name and type")
		}
		name, ok := args[0].(form.Term)
		if !ok {
			return nil, errors.New("define name must be a term")
		}
		typeExpr, err := ParseForm(args[1])
		if err != nil {
			return nil, fmt.Errorf("parsing define type: %w", err)
		}
		return Define{
			Name: Term(name),
			Type: typeExpr,
		}, nil

	case ":=":
		if len(args) != 2 {
			return nil, errors.New("assign must have exactly 2 arguments: name and value")
		}
		name, ok := args[0].(form.Term)
		if !ok {
			return nil, errors.New("assign name must be a term")
		}
		value, err := ParseForm(args[1])
		if err != nil {
			return nil, fmt.Errorf("parsing assign value: %w", err)
		}
		return Assign{
			Name:  Term(name),
			Value: value,
		}, nil

	case "match":
		if len(args) < 3 {
			return nil, errors.New("match must have at least 3 arguments: condition, cases, and final")
		}
		// First argument is the condition
		cond, err := ParseForm(args[0])
		if err != nil {
			return nil, fmt.Errorf("parsing match condition: %w", err)
		}
		// Parse cases - they come in pairs (comp, value)
		cases := make([]Case, 0)
		remainingArgs := args[1:]
		// Check if we have an even number of case arguments (pairs)
		// If odd, the last one is the final/default case
		var final Expr
		if len(remainingArgs)%2 != 1 {
			return nil, errors.New("match must have a final case")
		}

		// Last argument is the final/default case
		final, err = ParseForm(remainingArgs[len(remainingArgs)-1])
		if err != nil {
			return nil, fmt.Errorf("parsing match final case: %w", err)
		}
		remainingArgs = remainingArgs[:len(remainingArgs)-1]

		// Parse case pairs
		for i := 0; i < len(remainingArgs); i += 2 {
			comp, err := ParseForm(remainingArgs[i])
			if err != nil {
				return nil, fmt.Errorf("parsing match case %d comparison: %w", i/2, err)
			}
			value, err := ParseForm(remainingArgs[i+1])
			if err != nil {
				return nil, fmt.Errorf("parsing match case %d value: %w", i/2, err)
			}
			cases = append(cases, Case{
				Comp:  comp,
				Value: value,
			})
		}
		return Match{
			Cond:  cond,
			Cases: cases,
			Final: final,
		}, nil

	default:
		// Default to function call
		cmdExpr, err := ParseForm(cmd)
		if err != nil {
			return nil, fmt.Errorf("parsing function command: %w", err)
		}
		argExprs := make([]Expr, len(args))
		for i, arg := range args {
			argExpr, err := ParseForm(arg)
			if err != nil {
				return nil, fmt.Errorf("parsing function argument %d: %w", i, err)
			}
			argExprs[i] = argExpr
		}
		return FunctionCall{
			Cmd:  cmdExpr,
			Args: argExprs,
		}, nil
	}
}
