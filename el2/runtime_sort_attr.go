package el2

import "github.com/fbundle/sorts/sorts"

func (u Runtime) Form(s any) sorts.Form {
	return sorts.GetForm(u, s)
}

func (u Runtime) Level(s sorts.Sort) int {
	return sorts.GetLevel(u, s)
}
func (u Runtime) Parent(s sorts.Sort) sorts.Sort {
	return sorts.GetParent(u, s)
}
func (u Runtime) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return sorts.GetLessEqual(u, x, y)
}
func (u Runtime) TermOf(x sorts.Sort, X sorts.Sort) bool {
	return u.LessEqual(u.Parent(x), X)
}

func (u Runtime) NameLessEqual(src sorts.Name, dst sorts.Name) bool {
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
