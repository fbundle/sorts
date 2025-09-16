package form_processor2

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

func NewTokenizer(splitTokens []string) Tokenizer {
	s := slices.Clone(splitTokens)

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
		sortedSplitTokens: s,
	}
}

type Tokenizer struct {
	sortedSplitTokens []string
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

func (t Tokenizer) Tokenize(line string) (indentation int, fields []string) {
	indentation = 0
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
	fields = nil

	for _, field := range unSplitfields {
		for len(field) > 0 {
			if i, tok, ok := t.matchTok(field); ok {
				fmt.Printf("field \"%s\", match \"%s\" at %d\n", field, tok, i)
				beg, end := i, i+len(tok)
				fields, field = append(fields, field[:beg], field[beg:end]), field[end:]
				continue
			}
			fmt.Printf("field \"%s\", match nothing\n", field)
			fields, field = append(fields, field), ""

		}
	}

	return indentation, fields
}
