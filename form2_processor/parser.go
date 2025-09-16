package form2_processor

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

	if len(lines) == 0 {
		for range indentStack {
			code = append(code, p.CloseBlockToken)
		}
		return code
	}

	if indentStack[len(indentStack)-1] != lines[0].Indentation {
		panic("unreachable")
	}

	if len(lines) == 1 {
		for _, field := range lines[0].Fields {
			code = append(code, field)
		}
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
		for _, field := range lines[0].Fields {
			code = append(code, field)
		}
		code = append(code, p.NewLineToken)
		return p.parse(
			code,
			indentStack,
			lines[1:],
		)

	case nextInd > currInd: // open new block - do not add newline
		for _, field := range lines[0].Fields {
			code = append(code, field)
		}
		indentStack = append(indentStack, nextInd)
		return p.parse(
			code,
			indentStack,
			lines[1:],
		)
	case nextInd < currInd: // close block - add close block
		for nextInd < indentStack[len(indentStack)-1] {
			code = append(code, p.CloseBlockToken)
			indentStack = indentStack[:len(indentStack)-1]
		}
		return p.parse(
			code,
			indentStack,
			lines[1:],
		)
	default:
		panic("unreachable")
	}
}
