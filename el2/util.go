package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func mustSort(s almost_sort_extra.AlmostSort) almost_sort_extra.ActualSort {
	s1, ok := s.(almost_sort_extra.ActualSort)
	if !ok {
		panic(TypeErr)
	}
	return s1
}
func sortCompilerToAlmostSortCompiler(f func(form.Name) sorts.ListCompileFunc) func(form.Name) almost_sort_extra.ListCompileFunc {
	return func(name form.Name) almost_sort_extra.ListCompileFunc {
		sortCompiler := f(name)
		return func(r almost_sort_extra.Context, list form.List) (almost_sort_extra.Context, almost_sort_extra.AlmostSort) {
			return r, almost_sort_extra.NewActualSort(sortCompiler(func(form sorts.Form) sorts.Sort {
				var as almost_sort_extra.AlmostSort
				r, as = r.Compile(form)
				return mustSort(as).Repr()
			}, list))
		}
	}
}
