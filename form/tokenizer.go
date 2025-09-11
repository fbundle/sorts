package form

import (
	"sort"
	"strings"
	"unicode"
)

// getSplitTokens returns the sorted list of split tokens
func (parser Parser) getSplitTokens() []Token {
	splitTokenSet := make(map[Token]struct{})
	for _, tok := range parser.Split {
		splitTokenSet[tok] = struct{}{}
	}
	for beg, block := range parser.Blocks {
		splitTokenSet[beg] = struct{}{}
		splitTokenSet[block.End] = struct{}{}
	}
	splitTokens := make([]Token, 0, len(splitTokenSet))
	for tok := range splitTokenSet {
		splitTokens = append(splitTokens, tok)
	}
	sort.Slice(splitTokens, func(i, j int) bool {
		s1, s2 := splitTokens[i], splitTokens[j]
		if len(s1) > len(s2) && strings.HasPrefix(s1, s2) {
			// s2 is a prefix of s1, s1 should come first
			return true
		}
		if len(s2) > len(s1) && strings.HasPrefix(s2, s1) {
			// s1 is a prefix of s2, s2 should come first
			return false
		}
		// Otherwise, normal lexicographic order
		return s1 < s2
	})
	return splitTokens
}

type preprocessor func(string) string

var replaceAll = func(stringMap map[string]string) preprocessor {
	// replace all keys by the corresponding values
	return func(str string) string {
		for k, v := range stringMap {
			str = strings.ReplaceAll(str, k, v)
		}
		return str
	}
}

var removeComment = func(sep string) preprocessor {
	// drop content after sep in every line
	return func(str string) string {
		lines := strings.Split(str, "\n")
		var newLines []string
		for _, line := range lines {
			newLines = append(newLines, strings.SplitN(line, sep, 2)[0])
		}
		return strings.Join(newLines, "\n")
	}
}

const (
	CharStringBegin  Token = "\""
	CharStringEnd    Token = "\""
	CharStringEscape Token = "\\"
)

func tokenize(str string, splitTokens []string, pList ...preprocessor) []Token {
	// preprocess
	for _, p := range pList {
		str = p(str)
	}

	// state machine
	const (
		STATE_NORMAL = iota
		STATE_STRING
	)

	var tokens []Token
	state := STATE_NORMAL
	var buffer []rune
	flushBuffer := func() {
		if len(buffer) > 0 {
			tokens = append(tokens, string(buffer))
			buffer = buffer[:0]
		}
	}
	appendBuffer := func(n int) {
		str, buffer = str[n:], append(buffer, []rune(str[:n])...)
	}
	discardInput := func(n int) {
		str = str[n:]
	}

	for len(str) > 0 {
		switch state {
		case STATE_NORMAL:
			func() {
				for _, tok := range splitTokens {
					if len(str) >= len(tok) && str[:len(tok)] == tok {
						flushBuffer()
						appendBuffer(len(tok))
						flushBuffer()
						return
					}
				}
				if strings.HasPrefix(str, CharStringBegin) {
					flushBuffer()
					appendBuffer(len(CharStringBegin))
					state = STATE_STRING
				} else if unicode.IsSpace([]rune(str)[0]) {
					flushBuffer()
					discardInput(1)
				} else {
					appendBuffer(1)
				}
			}()
		case STATE_STRING:
			if len(str) > len(CharStringEscape) && strings.HasPrefix(str, CharStringEscape) {
				appendBuffer(len(CharStringEscape) + 1)
			} else if strings.HasPrefix(str, CharStringEnd) {
				appendBuffer(1)
				flushBuffer()
				state = STATE_NORMAL
			} else {
				appendBuffer(1)
			}
		default:
			panic("unreachable")
		}
	}

	flushBuffer()
	return tokens
}
