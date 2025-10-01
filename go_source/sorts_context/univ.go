package sorts_context

import (
	"cmp"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

// Univ - Initial, Terminal, Default must be distinct
type Univ struct {
	InitialTypeName  Name
	TerminalTypeName Name
	DefaultTypeName  Name
	nameLessEqual    ordered_map.Map[rule] // use Map since rule is not of cmp.Ordered
}

func (u Univ) Default(level int) Sort {
	return u.defaultSort(u.DefaultTypeName, level)
}
func (u Univ) Initial(level int) Sort {
	return u.defaultSort(u.InitialTypeName, level)
}
func (u Univ) Terminal(level int) Sort {
	return u.defaultSort(u.TerminalTypeName, level)
}

func (u Univ) defaultSort(name Name, level int) Sort {
	return sorts.NewTerm(
		setBuiltinLevel(name, level),
		sorts.NewChain(u.DefaultTypeName, level+1),
	)
}

func (u Univ) builtinNameGet(key Name) (Sort, bool) {
	// parse builtin: initial, terminal
	builtin := []Name{
		u.InitialTypeName,
		u.TerminalTypeName,
		u.DefaultTypeName,
	}
	name := string(key)
	for _, header := range builtin {
		if strings.HasPrefix(name, string(header)+"_") {
			levelStr := strings.TrimPrefix(name, string(header)+"_")
			level, err := strconv.Atoi(levelStr)
			if err != nil {
				continue
			}
			return u.defaultSort(header, level), true
		}
	}
	return nil, false
}

func (u Univ) LessEqual(s Form, d Form) bool {
	if equalForm(s, d) {
		return true
	}
	sName, ok1 := s.(Name)
	dName, ok2 := d.(Name)
	if ok1 {
		if _, ok := getBuiltinLevel(u.InitialTypeName, sName); ok {
			return true
		}
	}
	if ok2 {
		if _, ok := getBuiltinLevel(u.TerminalTypeName, dName); ok {
			return true
		}
	}
	if ok1 && ok2 {
		if _, ok := u.nameLessEqual.Get(rule{sName, dName}); ok {
			return true
		}
	}
	return false
}

func (u Univ) WithRule(src Name, dst Name) Univ {
	u.nameLessEqual = u.nameLessEqual.Set(rule{src, dst})
	return u
}

type rule struct {
	src Name
	dst Name
}

func (r rule) Cmp(s rule) int {
	if c := cmp.Compare(r.src, s.src); c != 0 {
		return c
	}
	return cmp.Compare(r.dst, s.dst)
}
