package sorts

import (
	"errors"
	"strconv"
	"strings"
)

var TypeErr = errors.New("type_error") // cannot recover

// Form - Union[Name, List]
type Form interface {
	mustForm()
}

type Name string

func (n Name) mustForm() {}

type List []Form

func (l List) mustForm() {}

type sortAttr struct {
	form      Form                // every Sort is identified with a Form
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

type Sort interface {
	sortAttr(a SortAttr) sortAttr
}

type SortAttr interface {
	Form(s any) Form
	Level(s Sort) int
	Parent(s Sort) Sort
	LessEqual(x Sort, y Sort) bool

	nameLessEqual(src Name, dst Name) bool
}

type Universe interface {
	Universe(level int) Atom
	Initial(level int) Atom
	Terminal(level int) Atom
	NewTerm(name Name, parent Sort) Atom

	NewNameLessEqualRule(src Name, dst Name)
	NewParseListRule(head Name, parseList ParseListFunc) error

	SortAttr
}

func NewUniverse(universeHeader Name, initialHeader Name, terminalHeader Name) (Universe, error) {
	nameSet := make(map[Name]struct{})
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

type ParseFunc = func(Form) (Sort, error)

type ParseListFunc = func(ParseFunc, List) (Sort, error)

type universe struct {
	universeHeader Name
	initialHeader  Name
	terminalHeader Name

	nameLessEqualDict map[[2]Name]struct{}
	parseListDict     map[Name]ParseListFunc

	nameDict map[Name]Atom
}

func (u *universe) Universe(level int) Atom {
	return newAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.universeHeader + "_" + levelStr
	})
}

func (u *universe) Initial(level int) Atom {
	levelStr := Name(strconv.Itoa(level))
	name := u.initialHeader + "_" + levelStr
	return newAtomTerm(u, name, u.Universe(level+1))
}

func (u *universe) Terminal(level int) Atom {
	levelStr := Name(strconv.Itoa(level))
	name := u.terminalHeader + "_" + levelStr
	return newAtomTerm(u, name, u.Universe(level+1))
}

func (u *universe) Parse(node Form) (Sort, error) {
	switch node := node.(type) {
	case Name:
		// lookup name
		if sort, ok := u.nameDict[node]; ok {
			return sort, nil
		}
		// parse builtin: universe, initial, terminal
		builtin := map[Name]func(level int) Atom{
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
	case List:
		if len(node) == 0 {
			return nil, errors.New("empty list")
		}
		head, ok := node[0].(Name)
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

func (u *universe) NewParseListRule(head Name, parseList ParseListFunc) error {
	if _, ok := u.parseListDict[head]; ok {
		return errors.New("list type already registered")
	}
	u.parseListDict[head] = parseList
	return nil
}

func (u *universe) NewNameLessEqualRule(src Name, dst Name) {
	u.nameLessEqualDict[[2]Name{src, dst}] = struct{}{}
}

func (u *universe) NewTerm(name Name, parent Sort) Atom {
	atom := newAtomTerm(u, name, parent)
	u.nameDict[name] = atom
	return atom
}

func (u *universe) Form(s any) Form {
	if sort, ok := s.(Sort); ok {
		return sort.sortAttr().form
	}
	if dep, ok := s.(Dependent); ok {
		return dep.Repr
	}
	panic(TypeErr)
}

func (u *universe) Level(s Sort) int {
	return s.sortAttr().level
}
func (u *universe) Parent(s Sort) Sort {
	return s.sortAttr().parent
}
func (u *universe) LessEqual(x Sort, y Sort) bool {
	return x.sortAttr().lessEqual(u, y)
}
func (u *universe) TermOf(x Sort, X Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u *universe) nameLessEqual(src Name, dst Name) bool {
	if src == u.initialHeader || dst == u.terminalHeader {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := u.nameLessEqualDict[[2]Name{src, dst}]; ok {
		return true
	}
	return false
}
