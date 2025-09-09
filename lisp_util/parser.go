package lisp_util

import "errors"

const (
	TokenBlockBegin Token = "("
	TokenBlockEnd   Token = ")"
	TokenSugarBegin Token = "{"
	TokenSugarEnd   Token = "}"
	TokenUnwrap     Token = "$"
	TokenTypeCast   Token = ":"
)

const (
	TermTypeCast     Term  = "type_cast"
	TermLambda       Term  = "lambda"
	SugarArrowSingle Token = "->"
	SugarArrowDouble Token = "=>"
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
	case TokenSugarBegin:
		// parse until seeing `}`
		argList, tokenList, err := parseUntilPred(Parse, matchName(Term(TokenSugarEnd)), tokenList)
		if err != nil {
			return nil, tokenList, err
		}
		expr, err := processSugar(argList)
		return expr, tokenList, err
	default:
		return Term(head), tokenList, nil
	}
}

// processSugar - handles both arithmetic infix and lambda syntax
// {1 + 2 + 3} -> (+ (+ 1 2) 3)
// {x y => (add x y)} -> (lambda x y (add x y))
// {x : type1} -> (type_cast type1 x)
func processSugar(argList []Expr) (Expr, error) {
	if len(argList) == 0 {
		return Node(nil), nil
	}
	if len(argList) == 1 {
		return argList[0], nil
	}
	secondLastName, ok := argList[len(argList)-2].(Term)
	if ok && string(secondLastName) == SugarArrowDouble {
		// arrow function syntax: {x y => expr}
		paramList := argList[:len(argList)-2]
		body := argList[len(argList)-1]
		lambdaArgList := []Expr{TermLambda}
		lambdaArgList = append(lambdaArgList, paramList...)
		lambdaArgList = append(lambdaArgList, body)

		return Node(lambdaArgList), nil
	}
	if ok && string(secondLastName) == TokenTypeCast {
		// type cast syntax
		typeCastArgList := []Expr{
			TermTypeCast,
			argList[len(argList)-1],
		}
		typeCastArgList = append(typeCastArgList, argList[:len(argList)-2]...)
		return Node(typeCastArgList), nil
	}

	// No arrow function or type cast, process as regular infix
	if ok && string(secondLastName) == SugarArrowSingle {
		// right to left
		left, cmd, argList := argList[0], argList[1], argList[1:]
		right, err := processSugar(argList)
		if err != nil {
			return nil, err
		}
		return Node([]Expr{cmd, left, right}), nil
	} else {
		// left to right
		argList, cmd, right := argList[:len(argList)-2], argList[len(argList)-2], argList[len(argList)-1]
		left, err := processSugar(argList)
		if err != nil {
			return nil, err
		}
		return Node([]Expr{cmd, left, right}), nil
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
