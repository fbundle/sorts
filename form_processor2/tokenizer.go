package form_processor2

import (
	"sort"
	"strings"
	"unicode"
)

type Line struct {
	Indentation int
	Fields      []string
}

type Tokenizer struct {
	SplitTokens       []string
	sortedSplitTokens []string
}

func (t Tokenizer) Init() Tokenizer {
	if len(t.SplitTokens) == len(t.sortedSplitTokens) {
		return t
	}

	s := t.SplitTokens
	for _, tok := range s {
		if len(tok) == 0 {
			panic("empty token")
		}
		for _, ch := range tok {
			if unicode.IsSpace(ch) {
				panic("tok cannot have space")
			}
		}
	}

	sort.Slice(s, func(i, j int) bool {
		s1, s2 := s[i], s[j]
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
	return Tokenizer{
		SplitTokens:       s,
		sortedSplitTokens: s,
	}
}

func (t Tokenizer) matchTok(s string) (int, string, bool) {
	matchIdx := len(s)
	matchTok := ""

	for _, tok := range t.sortedSplitTokens {
		i := strings.Index(s, tok)
		if i >= 0 {
			if i < matchIdx {
				matchIdx, matchTok = i, tok
			}
		}
	}
	return matchIdx, matchTok, matchIdx < len(s)
}

func (t Tokenizer) Tokenize(line string) Line {
	indentation := 0
	for _, r := range line {
		if unicode.IsSpace(r) && r != ' ' {
			panic("only space indentation allowed")
		}
		if r != ' ' {
			break
		}
		indentation++
	}

	unSplitfields := strings.Fields(line[indentation:])
	fields := make([]string, len(unSplitfields))

	consume := func(field string, length int) string {
		if length > 0 {
			fields = append(fields, field[:length])
		}
		return field[length:]
	}

	for _, field := range unSplitfields {
		for len(field) > 0 {
			if i, tok, ok := t.matchTok(field); ok {
				field = consume(field, i)
				field = consume(field, len(tok))
				continue
			} else {
				field = consume(field, len(field))
			}
		}
	}
	return Line{
		Indentation: indentation,
		Fields:      fields,
	}
}
