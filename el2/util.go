package el2

import (
	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func sortCompilerToAlmostSortCompiler(f func(form.Name) sorts.ListCompileFunc) func(form.Name) el_sorts.ListCompileFunc {
	return func(name form.Name) el_sorts.ListCompileFunc {
		compileFunc := f(name)
		return func(r el_sorts.Context, list form.List) el_sorts.Sort {
			return compileFunc(r.Compile, list)
		}
	}
}
