package el2

import (
	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func mustSort(s almost_sort_extra.Sort) almost_sort_extra.typeSort {
	s1, ok := s.(almost_sort_extra.typeSort)
	if !ok {
		panic(TypeErr)
	}
	return s1
}
func sortCompilerToAlmostSortCompiler(f func(form.Name) sorts.ListCompileFunc) func(form.Name) almost_sort_extra.ListCompileFunc {
	return func(name form.Name) almost_sort_extra.ListCompileFunc {
		sortCompiler := f(name)
		return func(r almost_sort_extra.Context, list form.List) (almost_sort_extra.Context, almost_sort_extra.Sort) {
			return r, almost_sort_extra.NewTypeSort(sortCompiler(func(form sorts.Form) sorts.Sort {
				var as almost_sort_extra.Sort
				r, as = r.Compile(form)
				return mustSort(as).Repr()
			}, list))
		}
	}
}
