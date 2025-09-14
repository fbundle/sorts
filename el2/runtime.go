package el2

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
)

func DefaultRuntime() Runtime {
	return Runtime{
		InitialHeader:  "Unit",
		TerminalHeader: "Any",
	}.
		NewListParser("->", toListParser(ListParseArrow("->"))).
		NewListParser("⊕", toListParser(ListParseSum("⊕"))).
		NewListParser("⊗", toListParser(ListParseProd("⊗"))).
		mustOk()
}

type Runtime struct {
	InitialHeader  Name
	TerminalHeader Name

	nameLessEqual ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered
	listParsers   ordered_map.OrderedMap[Name, ListParseFunc]

	frame ordered_map.OrderedMap[Name, Sort]
}

func (u Runtime) mustOk() Runtime {
	if u.InitialHeader == u.TerminalHeader {
		panic(TypeErr)
	}
	return u
}

// Terminal - T_0 T_1 ... T_n
func (u Runtime) Terminal(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.TerminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u Runtime) Initial(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.InitialHeader + "_" + levelStr
	})
}

func (u Runtime) Parse(node Form) AlmostSort {
	switch node := node.(type) {
	case Name:
		// lookup name
		if sort, ok := u.frame.Get(node); ok {
			return ActualSort{sort}
		}
		// parse builtin: Runtime, initial, terminal
		builtin := map[Name]func(level int) Sort{
			u.InitialHeader:  u.Initial,
			u.TerminalHeader: u.Terminal,
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

		listParser, ok := u.listParsers.Get(head)
		if !ok {
			panic("list type not registered")
		}
		// parse list
		return listParser(u.Parse, node)
	default:
		panic("parse error")
	}
}

func (u Runtime) NewListParser(head Name, parseList ListParseFunc) Runtime {
	if _, ok := u.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	return u.update(func(u Runtime) Runtime {
		u.listParsers = u.listParsers.Set(head, parseList)
		return u
	})
}

func (u Runtime) NewNameLessEqualRule(src Name, dst Name) Runtime {
	return u.update(func(u Runtime) Runtime {
		u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
		return u
	})
}

func (u Runtime) NewTerm(name Name, parent Sort) (Runtime, Sort) {
	atom := NewAtomTerm(u, name, parent)
	if _, ok := u.frame.Get(name); ok {
		panic("name already registered")
	}
	return u.update(func(u Runtime) Runtime {
		u.frame = u.frame.Set(name, atom)
		return u
	}), atom
}

func (u Runtime) update(f func(Runtime) Runtime) Runtime {
	return f(u)
}
