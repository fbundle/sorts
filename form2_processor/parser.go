package form2_processor

import "fmt"

type Parser struct {
	OpenBlockTokens []string
	CloseBlockToken string
	NewLineToken    string
}

func (p Parser) Parse(lines []Line) []string {
	if len(lines) == 0 {
		return nil
	}
	return p.parse(nil, []int{lines[0].Indentation}, lines)
}

func (p Parser) parse(code []string, indentStack []int, lines []Line) []string {
	fmt.Printf("parse: indent %v line %v\n", indentStack, lines)

	if len(lines) == 0 {
		for range indentStack {
			code = append(code, p.CloseBlockToken)
		}
		return code
	}

	if indentStack[len(indentStack)-1] != lines[0].Indentation {
		panic("unreachable")
	}

	for _, field := range lines[0].Fields {
		code = append(code, field)
	}

	if len(lines) == 1 {
		for range indentStack {
			code = append(code, p.CloseBlockToken)
		}
		return code
	}
	// len(lines) >= 2

	currInd := indentStack[len(indentStack)-1]
	nextInd := lines[1].Indentation

	switch {
	case nextInd == currInd: // same block - add new line
		code = append(code, p.NewLineToken)
	case nextInd < currInd: // open new block - do not add newline
		indentStack = append(indentStack, nextInd)
	case nextInd > currInd: // close block - add close block
		for nextInd < indentStack[len(indentStack)-1] {
			code = append(code, p.CloseBlockToken)
			indentStack = indentStack[:len(indentStack)-1]
		}
	default:
		panic("unreachable")
	}
	return p.parse(code, indentStack, lines[1:])
}
