package sort_universe

import (
	"cmp"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/el2"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type SortUniverse struct {
	initialHeader  form.Name
	terminalHeader form.Name
	nameLessEqual  ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered

}

func (u SortUniverse) mustSortAttr() SortUniverse {
	if u.initialHeader == u.terminalHeader {
		panic(el2.TypeErr)
	}
	return u
}

func (u SortUniverse) parseConstant(key form.Name) (el2.Sort, bool) {
	// parse builtin: initial, terminal
	builtin := map[form.Name]func(level int) el2.Sort{
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

func (u SortUniverse) NewNameLessEqualRule(src form.Name, dst form.Name) SortUniverse {
	u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
	return u
}

func (u SortUniverse) newTerm(name form.Name, parent el2.Sort) el2.Sort {
	return el2.NewAtomTerm(u, name, parent)
}

// Terminal - T_0 T_1 ... T_n
func (u SortUniverse) Terminal(level int) el2.Sort {
	return el2.NewAtomChain(level, func(level int) form.Name {
		levelStr := form.Name(strconv.Itoa(level))
		return u.terminalHeader + "_" + levelStr
	})
}

// Initial - I_0 I_1 ... I_n
func (u SortUniverse) Initial(level int) el2.Sort {
	return el2.NewAtomChain(level, func(level int) form.Name {
		levelStr := form.Name(strconv.Itoa(level))
		return u.initialHeader + "_" + levelStr
	})
}

func (u SortUniverse) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u SortUniverse) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u SortUniverse) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u SortUniverse) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u SortUniverse) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u SortUniverse) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
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
