package expr

import (
	"errors"
	"fmt"
)

func pop(tokenList []Token) ([]Token, Token, error) {
	if len(tokenList) == 0 {
		return nil, "", errors.New("empty token list")
	}
	return tokenList[1:], tokenList[0], nil
}

type Parser = func(tokenList []Token) (Expr, []Token, error)

func parseUntilPred(parser Parser, stopPred func(Expr) bool, tokenList []Token) ([]Expr, []Token, error) {
	var arg Expr
	var err error
	var argList []Expr
	for {
		arg, tokenList, err = parser(tokenList)
		if err != nil {
			return nil, tokenList, err
		}
		if stopPred(arg) {
			break
		}
		argList = append(argList, arg)
	}
	return argList, tokenList, nil
}

func Parse(tokenList []Token) (Expr, []Token, error) {
	tokenList, head, err := pop(tokenList)
	if err != nil {
		return nil, tokenList, err
	}

	switch head {
	case TokenBlockBegin:
		// parse until seeing `)`
		argList, tokenList, err := parseUntilPred(Parse, matchName(Term(TokenBlockEnd)), tokenList)
		if err != nil {
			return nil, tokenList, err
		}
		return Node(argList), tokenList, nil
	case TokenInfixBegin:
		// parse until seeing `}`
		argList, tokenList, err := parseUntilPred(Parse, matchName(Term(TokenInfixEnd)), tokenList)
		if err != nil {
			return nil, tokenList, err
		}
		expr, err := processInfix(argList)
		return expr, tokenList, err
	default:
		return Term(head), tokenList, nil
	}
}

const (
	TermSum         Term = "+"
	TermProd        Term = "×"
	TermArrowDouble Term = "=>"
	TermArrowSingle Term = "->"
	TermColon       Term = ":"
	TermComma       Term = ","
)

// processInfix - handles both infix and
// {1 + 2 + 3} 				(+ (+ 1 2) 3)				// left to right - sum
// {1 × 2 × 3}				(× (× 1 2) 3)				// left to right - prod
// {x => y => (add x y)}	(=> x (=> y (add x y)))		// right to left - lambda
// {x -> y -> z}			(-> x (-> y z))				// right to left - arrow
// {x : type1}				(: type1 x)					// right to left - type_cast
// {a, b, c}				(, a (, b c))				// right to left - list
func processInfix(argList []Expr) (Expr, error) {
	if len(argList) == 0 {
		return Node(nil), nil
	}
	if len(argList) == 1 {
		return argList[0], nil
	}
	if len(argList)%2 == 0 {
		return nil, errors.New("infix syntax must have an odd number of arguments")
	}
	op, ok := argList[1].(Term)
	if !ok {
		return nil, errors.New("infix operator must be a term")
	}
	for i := 3; i < len(argList); i += 2 {
		op2, ok := argList[i].(Term)
		if !ok {
			return nil, errors.New("infix operator must be a term")
		}
		if op2 != op {
			return nil, fmt.Errorf("infix operator must be the same %s", string(op))
		}
	}

	switch op {
	case TermSum, TermProd:
		// left to right
		argList, cmd, right := argList[:len(argList)-2], argList[len(argList)-2], argList[len(argList)-1]
		left, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return Node([]Expr{cmd, left, right}), nil
	case TermArrowDouble, TermArrowSingle, TermColon, TermComma:
		// right to left
		left, cmd, argList := argList[0], argList[1], argList[2:]
		right, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return Node([]Expr{cmd, left, right}), nil
	default:
		return nil, fmt.Errorf("infix operator not supported %s", string(op))
	}
}

func matchName(cond Term) func(Expr) bool {
	return func(arg Expr) bool {
		if name, ok := arg.(Term); ok {
			return string(cond) == string(name)
		}
		return false
	}
}
