package form_processor

import (
	"errors"
	"fmt"
)

var Tokenize = defaultProcessor.Tokenize

var Parse = defaultProcessor.Parse

type Preprocessor func(string) string
type PostProcessor func([]Token) []Token
type Block struct {
	End     Token
	Process func([]Form) (Form, error)
}

type Processor struct {
	Blocks map[Token]Block
	Split  []Token
}

var defaultProcessor = Processor{
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

func (p Processor) Tokenize(s string) []Token {
	preProcessorList := []Preprocessor{
		removeComment("#"),
	}
	postProcessorList := []PostProcessor{}

	return tokenize(
		s, p.getSplitTokens(),
		preProcessorList,
		postProcessorList,
	)
}

func (p Processor) Parse(tokenList []Token) ([]Token, Form, error) {
	tokenList, head, err := pop(tokenList)
	if err != nil {
		return tokenList, nil, err
	}
	if block, ok := p.Blocks[head]; ok {
		var f Form
		var forms []Form
		for {
			tokenList, f, err = p.Parse(tokenList)
			if err != nil {
				return tokenList, List(forms), err
			}
			if term, ok := f.(Name); ok && Token(term) == block.End {
				break
			}
			forms = append(forms, f)
		}
		f, err = block.Process(forms)
		return tokenList, f, err
	} else {
		return tokenList, Name(head), nil
	}
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
