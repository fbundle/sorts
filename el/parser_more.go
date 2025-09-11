package el

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

func ParseForm(e form.Form) (Expr, error) {
	return defaultParser.Parse(e)
}

var defaultParser = Parser{
	ListParsers: map[form.Term]ListParseFunc{
		"=>":    parseLambda,
		":":     parseDefine,
		":=":    parseAssign,
		"chain": parseChain,
		"match": parseMatch,
	},
}

func parseLambda(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) != 2 {
		return nil, errors.New("lambda must have exactly 2 arguments: param and body")
	}
	param, ok := list[0].(form.Term)
	if !ok {
		return nil, errors.New("lambda parameter must be a term")
	}
	body, err := parse(list[1])
	if err != nil {
		return nil, err
	}
	return Lambda{
		Param: Term(param),
		Body:  body,
	}, nil
}
func parseDefine(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) != 2 {
		return nil, errors.New("define must have exactly 2 arguments: name and type")
	}
	name, ok := list[0].(form.Term)
	if !ok {
		return nil, errors.New("define name must be a term")
	}
	typeExpr, err := parse(list[1])
	if err != nil {
		return nil, err
	}
	return Define{
		Name: Term(name),
		Type: typeExpr,
	}, nil
}
func parseAssign(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) != 2 {
		return nil, errors.New("assign must have exactly 2 arguments: name and value")
	}
	name, ok := list[0].(form.Term)
	if !ok {
		return nil, errors.New("assign name must be a term")
	}
	value, err := parse(list[1])
	if err != nil {
		return nil, err
	}
	return Assign{
		Name:  Term(name),
		Value: value,
	}, nil
}
func parseChain(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) < 2 {
		return nil, errors.New("chain must have at least 2 arguments: init expressions and tail")
	}
	// All arguments except the last are init expressions
	init := make([]Expr, len(list)-1)
	for i, arg := range list[:len(list)-1] {
		initExpr, err := parse(arg)
		if err != nil {
			return nil, err
		}
		init[i] = initExpr
	}
	// Last argument is the tail
	tail, err := parse(list[len(list)-1])
	if err != nil {
		return nil, err
	}
	return Chain{
		Init: init,
		Tail: tail,
	}, nil

}
func parseMatch(parse ParseFunc, list form.List) (Expr, error) {
	if len(list) < 3 {
		return nil, errors.New("match must have at least 3 arguments: condition, cases, and final")
	}
	// First argument is the condition
	cond, err := parse(list[0])
	if err != nil {
		return nil, err
	}
	// Parse cases - they come in pairs (comp, value)
	cases := make([]Case, 0)
	remainingArgs := list[1:]
	// Check if we have an odd number of remaining arguments (final case)
	var final Expr
	if len(remainingArgs)%2 != 1 {
		return nil, errors.New("match must have a final case")
	}
	// Last argument is the final/default case
	final, err = parse(remainingArgs[len(remainingArgs)-1])
	if err != nil {
		return nil, err
	}
	remainingArgs = remainingArgs[:len(remainingArgs)-1]
	// Parse case pairs
	for i := 0; i < len(remainingArgs); i += 2 {
		comp, err := parse(remainingArgs[i])
		if err != nil {
			return nil, err
		}
		value, err := parse(remainingArgs[i+1])
		if err != nil {
			return nil, err
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
}
