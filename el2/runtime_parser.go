package el2

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
)

type runtimeParser struct {
	listParsers ordered_map.OrderedMap[Name, ListParseFunc]
}

func (u runtimeParser) parse(node Form) AlmostSort {
	switch node := node.(type) {
	case Name:
		// lookup name
		if sort, ok := u.frame.Get(node); ok {
			return ActualSort{sort}
		}
		// parse builtin: initial, terminal
		builtin := map[Name]func(level int) Sort{
			u.initialHeader:  u.Initial,
			u.terminalHeader: u.Terminal,
		}
		name := string(node)
		for header, makeFunc := range builtin {
			if strings.HasPrefix(name, string(header)+"_") {
				levelStr := strings.TrimPrefix(name, string(header)+"_")
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					continue
				}
				sort := makeFunc(level)
				return ActualSort{sort}
			}
		}
		panic("name not found")
	case List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(Name)
		if !ok {
			panic("list must start with a name")
		}

		if listParser, ok := u.listParsers.Get(head); ok {
			return listParser(u.parse, node)
		} else { // by default, parse as beta reduction (function call)
			return ListParseBeta(u.parse, node)
		}
	default:
		panic("parse error")
	}
}

func (u runtimeParser) newListParser(head Name, parseList ListParseFunc) runtimeParser {
	if _, ok := u.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	u.listParsers = u.listParsers.Set(head, parseList)
	return u
}
