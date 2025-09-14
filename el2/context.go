package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/el_sorts"
	"github.com/fbundle/sorts/el2/universe"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

var _ almost_sort_extra.Context = Context{}

type Context struct {
	frame               ordered_map.OrderedMap[form.Name, almost_sort_extra.typeSort]
	universe            universe.SortUniverse
	listCompiler        ordered_map.OrderedMap[form.Name, almost_sort_extra.ListCompileFunc]
	defaultListCompiler almost_sort_extra.ListCompileFunc
}

func (ctx Context) Reset() Context {
	return Context{
		frame: ordered_map.OrderedMap[form.Name, almost_sort_extra.typeSort]{},
		universe: universe.SortUniverse{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
		},
		listCompiler:        ordered_map.OrderedMap[form.Name, almost_sort_extra.ListCompileFunc]{},
		defaultListCompiler: almost_sort_extra.ListCompileBeta,
	}.
		WithListCompiler("->", sortCompilerToAlmostSortCompiler(sorts.ListCompileArrow)).
		WithListCompiler("⊕", sortCompilerToAlmostSortCompiler(sorts.ListCompileSum)).
		WithListCompiler("⊗", sortCompilerToAlmostSortCompiler(sorts.ListCompileProd)).
		WithListCompiler("=>", almost_sort_extra.ListCompileLambda).
		WithListCompiler("let", almost_sort_extra.ListCompileLet("undef")).
		WithListCompiler("match", almost_sort_extra.ListCompileMatch("exact")).
		finalize()
}

// finalize - just for reset syntax to be neat
func (ctx Context) finalize() Context {
	return ctx
}
