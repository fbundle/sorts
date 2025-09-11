package form

import (
	"errors"
	"fmt"
)

type BlockConfig struct {
	End     Token
	Process func([]Form) (Form, error)
}

type Parser struct {
	Blocks map[Token]BlockConfig
	Split  []Token
}

var defaultParser = Parser{
	Split: []Token{"$", "⊕", "⊗", "=>", "->", ":", ",", "=", ":="},
	Blocks: map[Token]BlockConfig{
		"(": {
			End: ")",
			Process: func(forms []Form) (Form, error) {
				return List(forms), nil
			},
		},
		"{": {
			End:     "}",
			Process: processInfix,
		},
	},
}

func (parser Parser) Tokenize(s string) []Token {
	return tokenize(
		s, parser.getSplitTokens(),
		removeComment("#"),
	)
}

func (parser Parser) Parse(tokenList []Token) (Form, []Token, error) {
	tokenList, head, err := pop(tokenList)
	if err != nil {
		return nil, tokenList, err
	}
	if block, ok := parser.Blocks[head]; ok {
		var form Form
		var formList []Form
		for {
			form, tokenList, err = parser.Parse(tokenList)
			if err != nil {
				return List(formList), tokenList, err
			}
			if term, ok := form.(Term); ok && Token(term) == block.End {
				break
			}
			formList = append(formList, form)
		}
		form, err = block.Process(formList)
		return form, tokenList, err
	} else {
		return Term(head), tokenList, nil
	}
}

func Tokenize(s string) []Token {
	return defaultParser.Tokenize(s)
}

func Parse(tokenList []Token) (Form, []Token, error) {
	return defaultParser.Parse(tokenList)
}

func pop(tokenList []Token) ([]Token, Token, error) {
	if len(tokenList) == 0 {
		return nil, "", errors.New("empty token list")
	}
	return tokenList[1:], tokenList[0], nil
}

const (
	TermSum         Term = "⊕"
	TermProd        Term = "⊗"
	TermArrowDouble Term = "=>"
	TermArrowSingle Term = "->"
	TermColon       Term = ":"
	TermComma       Term = ","
	TermEqual       Term = "="
	TermColonEqual  Term = ":="
)

var leftToRightInfixOp = map[Term]struct{}{
	TermSum:  {}, // sum
	TermProd: {}, // prod
}

var rightToLeftInfixOp = map[Term]struct{}{
	TermArrowDouble: {}, // lambda expression
	TermArrowSingle: {}, // arrow type
	TermColon:       {}, // type cast
	TermComma:       {}, // list
	TermEqual:       {}, // equality
	TermColonEqual:  {}, // name binding
}

// processInfix - handles both infix and
// {1 + 2 + 3} 				(+ (+ 1 2) 3)				// left to right - sum
// {1 × 2 × 3}				(× (× 1 2) 3)				// left to right - prod
// {x => y => (add x y)}	(=> x (=> y (add x y)))		// right to left - lambda
// {x -> y -> z}			(-> x (-> y z))				// right to left - arrow
// etc.
func processInfix(argList []Form) (Form, error) {
	if len(argList) == 0 {
		return List(nil), nil
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

	if _, ok := leftToRightInfixOp[op]; ok {
		// left to right
		argList, cmd, right := argList[:len(argList)-2], argList[len(argList)-2], argList[len(argList)-1]
		left, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return List([]Form{cmd, left, right}), nil
	}
	if _, ok := rightToLeftInfixOp[op]; ok {
		// right to left
		left, cmd, argList := argList[0], argList[1], argList[2:]
		right, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return List([]Form{cmd, left, right}), nil
	}

	return nil, fmt.Errorf("infix operator not supported %s", string(op))
}
