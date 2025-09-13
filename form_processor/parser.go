package form_processor

import (
	"errors"
	"fmt"
)

type Preprocessor func(string) string
type Block struct {
	End     Token
	Process func([]Form) (Form, error)
}

type Parser struct {
	Blocks map[Token]Block
	Split  []Token
}

var defaultParser = Parser{
	Split: []Token{"+", "*", "$", "⊕", "⊗", "Π", "Σ", "=>", "->", ":", ",", "=", ":="},
	Blocks: map[Token]Block{
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

func (parser Parser) Tokenize(s string, pList ...Preprocessor) []Token {
	newPList := append([]Preprocessor{
		removeComment("#"),
	}, pList...)

	return tokenize(
		s, parser.getSplitTokens(),
		newPList...,
	)
}

func (parser Parser) Parse(tokenList []Token) (Form, []Token, error) {
	tokenList, head, err := pop(tokenList)
	if err != nil {
		return nil, tokenList, err
	}
	if block, ok := parser.Blocks[head]; ok {
		var f Form
		var forms []Form
		for {
			f, tokenList, err = parser.Parse(tokenList)
			if err != nil {
				return List(forms), tokenList, err
			}
			if term, ok := f.(Name); ok && Token(term) == block.End {
				break
			}
			forms = append(forms, f)
		}
		f, err = block.Process(forms)
		return f, tokenList, err
	} else {
		return Name(head), tokenList, nil
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

// processInfix - handles both infix and
// {1 + 2 + 3} 				(+ (+ 1 2) 3)				// left assoc
// {1 × 2 × 3}				(× (× 1 2) 3)				// left assoc
// {x => y => (add x y)}	(=> x (=> y (add x y)))		// right assoc
// {x -> y -> z}			(-> x (-> y z))				// right assoc
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
	op, ok := argList[1].(Name)
	if !ok {
		return nil, errors.New("infix operator must be a term")
	}
	for i := 3; i < len(argList); i += 2 {
		op2, ok := argList[i].(Name)
		if !ok {
			return nil, errors.New("infix operator must be a term")
		}
		if op2 != op {
			return nil, fmt.Errorf("infix operator must be the same %s", string(op))
		}
	}

	leftAssocOperator := map[Name]struct{}{
		"+": {},
		"*": {},
	}

	if _, ok := leftAssocOperator[op]; ok {
		// left assoc
		argList, cmd, right := argList[:len(argList)-2], argList[len(argList)-2], argList[len(argList)-1]
		left, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return List([]Form{cmd, left, right}), nil
	} else {
		// by default, right assoc
		left, cmd, argList := argList[0], argList[1], argList[2:]
		right, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return List([]Form{cmd, left, right}), nil
	}
}
