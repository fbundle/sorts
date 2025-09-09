package ast

import (
	"strings"
)

const (
	TokenBlockBegin Token = "("
	TokenBlockEnd   Token = ")"
	TokenInfixBegin Token = "{"
	TokenInfixEnd   Token = "}"
	TokenUnwrap     Token = "$"
	TokenSum        Token = "+"
	TokenProd       Token = "×"
	TokenLambda     Token = "=>"
	TokenTypeCast   Token = ":"
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
		},
		removeComment("#"),
		replaceAll(map[string]string{
			"[": " (list ",
			"]": " ) ",
		}),
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
	TokenStringBegin Token = "\""
	TokenStringEnd   Token = "\""
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

	for len(str) > 0 {
		switch state {
		case STATE_NORMAL:
		case STATE_STRING:
		default:
			panic("unreachable")
		}
	}

	flushBuffer()
	return tokens
}
