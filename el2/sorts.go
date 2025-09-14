package el2

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type sortUniverse struct {
	initialHeader  Name
	terminalHeader Name
	nameLessEqual  ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered

}

func (u sortUniverse) mustSortAttr() sortUniverse {
	if u.initialHeader == u.terminalHeader {
		panic(TypeErr)
	}
	return u
}

func (u sortUniverse) parseConstant(key Name) (Sort, bool) {
	// parse builtin: initial, terminal
	builtin := map[Name]func(level int) Sort{
		u.initialHeader:  u.Initial,
		u.terminalHeader: u.Terminal,
	}
	name := string(key)
	for header, makeFunc := range builtin {
		if strings.HasPrefix(name, string(header)+"_") {
			levelStr := strings.TrimPrefix(name, string(header)+"_")
			level, err := strconv.Atoi(levelStr)
			if err != nil {
				continue
			}
			sort := makeFunc(level)
			return sort, true
		}
	}
	return nil, false
}

func (u sortUniverse) NewNameLessEqualRule(src Name, dst Name) sortUniverse {
	u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
	return u
}

func (u sortUniverse) newTerm(name Name, parent Sort) Sort {
	return NewAtomTerm(u, name, parent)
}

// Terminal - T_0 T_1 ... T_n
func (u sortUniverse) Terminal(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.terminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u sortUniverse) Initial(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.initialHeader + "_" + levelStr
	})
}

func (u sortUniverse) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u sortUniverse) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u sortUniverse) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u sortUniverse) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u sortUniverse) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u sortUniverse) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
	if src == u.initialHeader || dst == u.terminalHeader {
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
