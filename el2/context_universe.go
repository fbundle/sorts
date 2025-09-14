package el2

import "github.com/fbundle/sorts/sorts"

func (r Context) Form(s any) sorts.Form {
	return r.sortUniverse.Form(s)
}

func (r Context) Level(s sorts.Sort) int {
	return r.sortUniverse.Level(s)
}

func (r Context) Parent(s sorts.Sort) sorts.Sort {
	return r.sortUniverse.Parent(s)
}

func (r Context) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return r.sortUniverse.LessEqual(x, y)
}

func (r Context) GetRule(src sorts.Name, dst sorts.Name) bool {
	return r.sortUniverse.GetRule(src, dst)
}
