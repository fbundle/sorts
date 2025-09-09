package lisp_util

import (
	"fmt"
	"strings"
	"unicode"
)

type Char = rune

const (
	CharBlockBegin   Char = '('
	CharBlockEnd     Char = ')'
	CharSugarBegin   Char = '{'
	CharSugarEnd     Char = '}'
	CharUnwrap       Char = '$'
	CharTypeCast     Char = ':'
	CharStringBegin  Char = '"'
	CharStringEnd    Char = '"'
	CharStringEscape Char = '\\'
)

func Tokenize(s string) []Token {
	return tokenize(s,
		map[Char]struct{}{
			CharBlockBegin: {},
			CharBlockEnd:   {},
			CharSugarBegin: {},
			CharSugarEnd:   {},
			CharUnwrap:     {},
			CharTypeCast:   {},
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

func tokenize(str string, splitChars map[Char]struct{}, pList ...preprocessor) []Token {
	// preprocess
	for _, p := range pList {
		str = p(str)
	}

	// state machine
	const (
		STATE_NORMAL = iota
		STATE_STRING
		STATE_STRING_ESCAPE
	)

	var tokens []Token
	state := STATE_NORMAL
	var buffer []rune = nil
	flushBuffer := func() {
		if len(buffer) > 0 {
			tokens = append(tokens, string(buffer))
		}
		buffer = nil
	}
	appendBuffer := func(ch Char) {
		buffer = append(buffer, ch)
	}

	for _, ch := range str {
		switch state {
		case STATE_NORMAL: // outside string
			if _, ok := splitChars[ch]; ok {
				// split special characters like ( ) [ ] * : into tokens
				flushBuffer()
				appendBuffer(ch)
				flushBuffer()
			} else if unicode.IsSpace(ch) {
				// flush buffer if seeing whitespace
				flushBuffer()
			} else if ch == CharStringBegin {
				// enter string mode
				flushBuffer()
				appendBuffer(ch)
				state = STATE_STRING
			} else {
				appendBuffer(ch)
			}
		case STATE_STRING:
			if ch == CharStringEscape {
				appendBuffer(ch)
				state = STATE_STRING_ESCAPE
			} else if ch == CharStringEnd {
				// exit string mode
				appendBuffer(ch)
				flushBuffer()
				state = STATE_NORMAL
			} else {
				appendBuffer(ch)
			}
		case STATE_STRING_ESCAPE:
			appendBuffer(ch)
			state = STATE_STRING
		default:
			panic(fmt.Sprintf("unreachable state: %d", state))
		}
	}
	flushBuffer()
	return tokens
}
