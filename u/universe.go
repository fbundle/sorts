package u

import (
	"cmp"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Universe interface {
	sorts.SortAttr

	Initial(level int) sorts.Atom
	Terminal(level int) sorts.Atom
	NewTerm(name sorts.Name, parent sorts.Sort) (Universe, sorts.Atom)

	NewNameLessEqualRule(src sorts.Name, dst sorts.Name) Universe
	NewParseListRule(head sorts.Name, parseList sorts.ParseListFunc) Universe
}

func newDefaultUniverse() Universe {
	return newUniverse("Unit", "Any").
		NewParseListRule("->", sorts.ParseListArrow("->"))
}

func newUniverse(initialHeader sorts.Name, terminalHeader sorts.Name) Universe {
	nameSet := make(map[sorts.Name]struct{})
	nameSet[initialHeader] = struct{}{}
	nameSet[terminalHeader] = struct{}{}
	if len(nameSet) != 3 {
		panic("universe, initial, terminal name must be distinct")
	}
	u := &universe{
		initialHeader:  initialHeader,
		terminalHeader: terminalHeader,
	}
	return u
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

type universe struct {
	initialHeader  sorts.Name
	terminalHeader sorts.Name

	nameLessEqualDict ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered
	parseListDict     ordered_map.OrderedMap[sorts.Name, sorts.ParseListFunc]

	nameDict ordered_map.OrderedMap[sorts.Name, sorts.Atom]
}

// Terminal - T_0 T_1 ... T_n
func (u universe) Terminal(level int) sorts.Atom {
	return sorts.NewAtomChain(level, func(level int) sorts.Name {
		levelStr := sorts.Name(strconv.Itoa(level))
		return u.terminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u universe) Initial(level int) sorts.Atom {
	return sorts.NewAtomChain(level, func(level int) sorts.Name {
		levelStr := sorts.Name(strconv.Itoa(level))
		return u.initialHeader + "_" + levelStr
	})
}

func (u universe) Parse(node sorts.Form) sorts.Sort {
	switch node := node.(type) {
	case sorts.Name:
		// lookup name
		if sort, ok := u.nameDict.Get(node); ok {
			return sort
		}
		// parse builtin: universe, initial, terminal
		builtin := map[sorts.Name]func(level int) sorts.Atom{
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

		rule, ok := u.parseListDict.Get(head)
		if !ok {
			panic("list type not registered")
		}
		// parse list
		sort, err := rule(func(form sorts.Form) (sorts.Sort, error) {
			return u.Parse(form), nil
		}, node)
		if err != nil {
			panic(err)
		}
		return sort
	default:
		panic("parse error")
	}
}

func (u universe) NewParseListRule(head sorts.Name, parseList sorts.ParseListFunc) Universe {
	if _, ok := u.parseListDict.Get(head); ok {
		panic("list type already registered")
	}
	return u.update(func(u universe) universe {
		u.parseListDict = u.parseListDict.Set(head, parseList)
		return u
	})
}

func (u universe) NewNameLessEqualRule(src sorts.Name, dst sorts.Name) Universe {
	return u.update(func(u universe) universe {
		u.nameLessEqualDict = u.nameLessEqualDict.Set(rule{src, dst})
		return u
	})
}

func (u universe) NewTerm(name sorts.Name, parent sorts.Sort) (Universe, sorts.Atom) {
	atom := sorts.NewAtomTerm(u, name, parent)
	if _, ok := u.nameDict.Get(name); ok {
		panic("name already registered")
	}
	return u.update(func(u universe) universe {
		u.nameDict = u.nameDict.Set(name, atom)
		return u
	}), atom
}

func (u universe) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u universe) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u universe) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u universe) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u universe) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u universe) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
	if src == u.initialHeader || dst == u.terminalHeader {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := u.nameLessEqualDict.Get(rule{src, dst}); ok {
		return true
	}
	return false
}

func (u universe) update(f func(universe) universe) universe {
	return f(u)
}
