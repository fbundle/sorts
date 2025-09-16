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

	fmt.Println(s)

	return Tokenizer{
		sortedSplitTokens: s,
	}
}

type Tokenizer struct {
	sortedSplitTokens []string
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
			matchIdx := -1
			matchTok := ""
			for _, tok := range t.sortedSplitTokens {
				index := strings.Index(field, tok)
				if index >= 0 {
					matchIdx, matchTok = index, tok
					break
				}
			}
			if matchIdx < 0 {
				fields, field = append(fields, field), ""
			} else {
				beg, end := matchIdx, matchIdx+len(matchTok)
				fields, field = append(fields, field[:beg], field[beg:end]), field[end:]
			}
		}
	}

	return indentation, fields
}
