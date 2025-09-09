package ast

import (
	"strings"
	"unicode"
)

const (
	TokenBlockBegin Token = "("
	TokenBlockEnd   Token = ")"
	TokenInfixBegin Token = "{"
	TokenInfixEnd   Token = "}"
	TokenUnwrap     Token = "$"

	TokenSum      Token = "+"
	TokenProd     Token = "×"
	TokenLambda   Token = "=>"
	TokenTypeCast Token = ":"
	TokenList     Token = ","
)

func Tokenize(s string) []Token {
	return tokenize(s,
		[]string{
			TokenBlockBegin,
			TokenBlockEnd,
			TokenInfixBegin,
			TokenInfixEnd,
			TokenUnwrap,
			TokenSum,
			TokenProd,
			TokenLambda,
			TokenTypeCast,
			TokenList,
		},
		removeComment("#"),
	)
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
				if str[0:1] == CharStringBegin {
					flushBuffer()
					appendBuffer(1)
					state = STATE_STRING
				} else if unicode.IsSpace([]rune(str)[0]) {
					flushBuffer()
					discardInput(1)
				} else {
					appendBuffer(1)
				}
			}()
		case STATE_STRING:
			if len(str) >= 2 && str[0:1] == CharStringEscape {
				appendBuffer(2)
			} else if str[0:1] == CharStringEnd {
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
