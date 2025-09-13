package form_processor

import (
	"errors"
	"fmt"

	"github.com/fbundle/sorts/form"
)

type Preprocessor func(string) string
type Block struct {
	End     form.Token
	Process func([]form.Form) (form.Form, error)
}

type Parser struct {
	Blocks map[form.Token]Block
	Split  []form.Token
}

var defaultParser = Parser{
	Split: []form.Token{"+", "*", "$", "⊕", "⊗", "Π", "Σ", "=>", "->", ":", ",", "=", ":="},
	Blocks: map[form.Token]Block{
		"(": {
			End: ")",
			Process: func(forms []form.Form) (form.Form, error) {
				return form.List(forms), nil
			},
		},
		"{": {
			End:     "}",
			Process: processInfix,
		},
	},
}

func (parser Parser) Tokenize(s string, pList ...Preprocessor) []form.Token {
	newPList := append([]Preprocessor{
		removeComment("#"),
	}, pList...)

	return tokenize(
		s, parser.getSplitTokens(),
		newPList...,
	)
}

func (parser Parser) Parse(tokenList []form.Token) (form.Form, []form.Token, error) {
	tokenList, head, err := pop(tokenList)
	if err != nil {
		return nil, tokenList, err
	}
	if block, ok := parser.Blocks[head]; ok {
		var form form.Form
		var formList []form.Form
		for {
			form, tokenList, err = parser.Parse(tokenList)
			if err != nil {
				return form.List(formList), tokenList, err
			}
			if term, ok := form.(form.Name); ok && form.Token(term) == block.End {
				break
			}
			formList = append(formList, form)
		}
		form, err = block.Process(formList)
		return form, tokenList, err
	} else {
		return form.Name(head), tokenList, nil
	}
}

func Tokenize(s string) []form.Token {
	return defaultParser.Tokenize(s)
}

func Parse(tokenList []form.Token) (form.Form, []form.Token, error) {
	return defaultParser.Parse(tokenList)
}

func pop(tokenList []form.Token) ([]form.Token, form.Token, error) {
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
func processInfix(argList []form.Form) (form.Form, error) {
	if len(argList) == 0 {
		return form.List(nil), nil
	}
	if len(argList) == 1 {
		return argList[0], nil
	}
	if len(argList)%2 == 0 {
		return nil, errors.New("infix syntax must have an odd number of arguments")
	}
	op, ok := argList[1].(form.Name)
	if !ok {
		return nil, errors.New("infix operator must be a term")
	}
	for i := 3; i < len(argList); i += 2 {
		op2, ok := argList[i].(form.Name)
		if !ok {
			return nil, errors.New("infix operator must be a term")
		}
		if op2 != op {
			return nil, fmt.Errorf("infix operator must be the same %s", string(op))
		}
	}

	leftAssocOperator := map[form.Name]struct{}{
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
		return form.List([]form.Form{cmd, left, right}), nil
	} else {
		// by default, right assoc
		left, cmd, argList := argList[0], argList[1], argList[2:]
		right, err := processInfix(argList)
		if err != nil {
			return nil, err
		}
		return form.List([]form.Form{cmd, left, right}), nil
	}
}
