package el2

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type runtimeSortAttr struct {
	initialHeader  Name
	terminalHeader Name
	nameLessEqual  ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered

}

func (u runtimeSortAttr) mustSortAttr() runtimeSortAttr {
	if u.initialHeader == u.terminalHeader {
		panic(TypeErr)
	}
	return u
}

func (u runtimeSortAttr) parseConstant(key Name) (Sort, bool) {
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

func (u runtimeSortAttr) NewNameLessEqualRule(src Name, dst Name) runtimeSortAttr {
	u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
	return u
}

func (u runtimeSortAttr) newTerm(name Name, parent Sort) Sort {
	return NewAtomTerm(u, name, parent)
}

// Terminal - T_0 T_1 ... T_n
func (u runtimeSortAttr) Terminal(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.terminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u runtimeSortAttr) Initial(level int) Sort {
	return NewAtomChain(level, func(level int) Name {
		levelStr := Name(strconv.Itoa(level))
		return u.initialHeader + "_" + levelStr
	})
}

func (u runtimeSortAttr) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u runtimeSortAttr) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u runtimeSortAttr) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u runtimeSortAttr) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u runtimeSortAttr) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u runtimeSortAttr) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
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
