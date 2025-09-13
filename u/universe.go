package u

import (
	"cmp"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

func DefaultUniverse() Universe {
	return Universe{
		InitialHeader:  "Unit",
		TerminalHeader: "Any",
	}.
		NewListParser("->", sorts.ListParseArrow("->")).
		NewListParser("⊕", sorts.ListParseSum("⊕")).
		NewListParser("⊗", sorts.ListParseProd("⊗"))
}

type rule struct {
	src sorts.Name
	dst sorts.Name
}

func (r rule) Cmp(s rule) int {
	if c := cmp.Compare(r.src, s.src); c != 0 {
		return c
	}
	return cmp.Compare(r.dst, s.dst)
}

type Universe struct {
	InitialHeader  sorts.Name
	TerminalHeader sorts.Name

	nameLessEqual ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered
	listParsers   ordered_map.OrderedMap[sorts.Name, sorts.ListParseFunc]

	frame ordered_map.OrderedMap[sorts.Name, sorts.Sort]
}

// Terminal - T_0 T_1 ... T_n
func (u Universe) Terminal(level int) sorts.Sort {
	return sorts.NewAtomChain(level, func(level int) sorts.Name {
		levelStr := sorts.Name(strconv.Itoa(level))
		return u.TerminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u Universe) Initial(level int) sorts.Sort {
	return sorts.NewAtomChain(level, func(level int) sorts.Name {
		levelStr := sorts.Name(strconv.Itoa(level))
		return u.InitialHeader + "_" + levelStr
	})
}

func (u Universe) Parse(node sorts.Form) sorts.Sort {
	switch node := node.(type) {
	case sorts.Name:
		// lookup name
		if sort, ok := u.frame.Get(node); ok {
			return sort
		}
		// parse builtin: Universe, initial, terminal
		builtin := map[sorts.Name]func(level int) sorts.Sort{
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
				return sort
			}
		}
		panic("name not found")
	case sorts.List:
		if len(node) == 0 {
			panic("empty list")
		}
		head, ok := node[0].(sorts.Name)
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

func (u Universe) NewListParser(head sorts.Name, parseList sorts.ListParseFunc) Universe {
	if _, ok := u.listParsers.Get(head); ok {
		panic("list type already registered")
	}
	return u.update(func(u Universe) Universe {
		u.listParsers = u.listParsers.Set(head, parseList)
		return u
	})
}

func (u Universe) NewNameLessEqualRule(src sorts.Name, dst sorts.Name) Universe {
	return u.update(func(u Universe) Universe {
		u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
		return u
	})
}

func (u Universe) NewTerm(name sorts.Name, parent sorts.Sort) (Universe, sorts.Sort) {
	atom := sorts.NewAtomTerm(u, name, parent)
	if _, ok := u.frame.Get(name); ok {
		panic("name already registered")
	}
	return u.update(func(u Universe) Universe {
		u.frame = u.frame.Set(name, atom)
		return u
	}), atom
}

func (u Universe) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u Universe) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u Universe) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u Universe) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u Universe) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u Universe) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
	if src == u.InitialHeader || dst == u.TerminalHeader {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := u.nameLessEqual.Get(rule{src, dst}); ok {
		return true
	}
	return false
}

func (u Universe) update(f func(Universe) Universe) Universe {
	return f(u)
}
