package sort_universe

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/el2"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type sortUniverse struct {
	initialHeader  el2.Name
	terminalHeader el2.Name
	nameLessEqual  ordered_map.Map[el2.rule] // use Map since rule is not of cmp.Ordered

}

func (u sortUniverse) mustSortAttr() sortUniverse {
	if u.initialHeader == u.terminalHeader {
		panic(el2.TypeErr)
	}
	return u
}

func (u sortUniverse) parseConstant(key el2.Name) (el2.Sort, bool) {
	// parse builtin: initial, terminal
	builtin := map[el2.Name]func(level int) el2.Sort{
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

func (u sortUniverse) NewNameLessEqualRule(src el2.Name, dst el2.Name) sortUniverse {
	u.nameLessEqual = u.nameLessEqual.Set(el2.rule{src, dst})
	return u
}

func (u sortUniverse) newTerm(name el2.Name, parent el2.Sort) el2.Sort {
	return el2.NewAtomTerm(u, name, parent)
}

// Terminal - T_0 T_1 ... T_n
func (u sortUniverse) Terminal(level int) el2.Sort {
	return el2.NewAtomChain(level, func(level int) el2.Name {
		levelStr := el2.Name(strconv.Itoa(level))
		return u.terminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u sortUniverse) Initial(level int) el2.Sort {
	return el2.NewAtomChain(level, func(level int) el2.Name {
		levelStr := el2.Name(strconv.Itoa(level))
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
	if _, ok := u.nameLessEqual.Get(el2.rule{src, dst}); ok {
		return true
	}
	return false
}
