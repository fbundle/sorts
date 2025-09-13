package universe

import (
	"errors"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/sorts"
)

type Universe interface {
	sorts.SortAttr

	Universe(level int) sorts.Atom
	Initial(level int) sorts.Atom
	Terminal(level int) sorts.Atom
	NewTerm(name sorts.Name, parent sorts.Sort) sorts.Atom

	NewNameLessEqualRule(src sorts.Name, dst sorts.Name)
	NewParseListRule(head sorts.Name, parseList ParseListFunc) error
}

func NewUniverse(universeHeader sorts.Name, initialHeader sorts.Name, terminalHeader sorts.Name) (Universe, error) {
	nameSet := make(map[sorts.Name]struct{})
	nameSet[universeHeader] = struct{}{}
	nameSet[initialHeader] = struct{}{}
	nameSet[terminalHeader] = struct{}{}
	if len(nameSet) != 3 {
		return nil, errors.New("universe, initial, terminal name must be distinct")
	}
	return &universe{
		universeHeader: universeHeader,
		initialHeader:  initialHeader,
		terminalHeader: terminalHeader,
	}, nil
}

type ParseFunc = func(sorts.Form) (sorts.Sort, error)

type ParseListFunc = func(ParseFunc, sorts.List) (sorts.Sort, error)

type universe struct {
	universeHeader sorts.Name
	initialHeader  sorts.Name
	terminalHeader sorts.Name

	nameLessEqualDict map[[2]sorts.Name]struct{}
	parseListDict     map[sorts.Name]ParseListFunc

	nameDict map[sorts.Name]sorts.Atom
}

func (u *universe) Universe(level int) sorts.Atom {
	return sorts.NewAtomChain(level, func(level int) sorts.Name {
		levelStr := sorts.Name(strconv.Itoa(level))
		return u.universeHeader + "_" + levelStr
	})
}

func (u *universe) Initial(level int) sorts.Atom {
	levelStr := sorts.Name(strconv.Itoa(level))
	name := u.initialHeader + "_" + levelStr
	return sorts.NewAtomTerm(u, name, u.Universe(level+1))
}

func (u *universe) Terminal(level int) sorts.Atom {
	levelStr := sorts.Name(strconv.Itoa(level))
	name := u.terminalHeader + "_" + levelStr
	return sorts.NewAtomTerm(u, name, u.Universe(level+1))
}

func (u *universe) Parse(node sorts.Form) (sorts.Sort, error) {
	switch node := node.(type) {
	case sorts.Name:
		// lookup name
		if sort, ok := u.nameDict[node]; ok {
			return sort, nil
		}
		// parse builtin: universe, initial, terminal
		builtin := map[sorts.Name]func(level int) sorts.Atom{
			u.universeHeader: u.Universe,
			u.initialHeader:  u.Initial,
			u.terminalHeader: u.Terminal,
		}
		name := string(node)
		for header, makeFunc := range builtin {
			if strings.HasPrefix(name, string(header)+"_") {
				levelStr := strings.TrimPrefix(name, string(header)+"_")
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return nil, err
				}
				sort := makeFunc(level)
				return sort, nil
			}
		}
		return nil, errors.New("name not found")
	case sorts.List:
		if len(node) == 0 {
			return nil, errors.New("empty list")
		}
		head, ok := node[0].(sorts.Name)
		if !ok {
			return nil, errors.New("list must start with a name")
		}

		rule, ok := u.parseListDict[head]
		if !ok {
			return nil, errors.New("list type not registered")
		}
		// parse list
		return rule(u.Parse, node[1:])
	default:
		return nil, errors.New("parse error")
	}
}

func (u *universe) NewParseListRule(head sorts.Name, parseList ParseListFunc) error {
	if _, ok := u.parseListDict[head]; ok {
		return errors.New("list type already registered")
	}
	u.parseListDict[head] = parseList
	return nil
}

func (u *universe) NewNameLessEqualRule(src sorts.Name, dst sorts.Name) {
	u.nameLessEqualDict[[2]sorts.Name{src, dst}] = struct{}{}
}

func (u *universe) NewTerm(name sorts.Name, parent sorts.Sort) sorts.Atom {
	atom := sorts.NewAtomTerm(u, name, parent)
	u.nameDict[name] = atom
	return atom
}

func (u *universe) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u *universe) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u *universe) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u *universe) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u *universe) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u *universe) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
	if src == u.initialHeader || dst == u.terminalHeader {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := u.nameLessEqualDict[[2]sorts.Name{src, dst}]; ok {
		return true
	}
	return false
}
