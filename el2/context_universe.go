package el2

import "github.com/fbundle/sorts/sorts"

func (ctx Context) Form(s any) sorts.Form {
	return ctx.sortUniverse.Form(s)
}

func (ctx Context) Level(s sorts.Sort) int {
	return ctx.sortUniverse.Level(s)
}

func (ctx Context) Parent(s sorts.Sort) sorts.Sort {
	return ctx.sortUniverse.Parent(s)
}

func (ctx Context) LessEqual(x sorts.Sort, y sorts.Sort) bool {
	return ctx.sortUniverse.LessEqual(x, y)
}

func (ctx Context) GetRule(src sorts.Name, dst sorts.Name) bool {
	return ctx.sortUniverse.GetRule(src, dst)
}
