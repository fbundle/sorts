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

var _ el_sorts.Context = Context{}

type Context struct {
	frame               ordered_map.OrderedMap[form.Name, el_sorts.Sort]
	universe            universe.SortUniverse
	listCompiler        ordered_map.OrderedMap[form.Name, el_sorts.ListCompileFunc]
	defaultListCompiler el_sorts.ListCompileFunc
}

func (ctx Context) Reset() Context {
	return Context{
		frame: ordered_map.OrderedMap[form.Name, el_sorts.Sort]{},
		universe: universe.SortUniverse{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
		},
		listCompiler:        ordered_map.OrderedMap[form.Name, el_sorts.ListCompileFunc]{},
		defaultListCompiler: el_sorts.ListCompileBeta,
	}.
		WithListCompiler("->", sortCompilerToAlmostSortCompiler(sorts.ListCompileArrow)).
		WithListCompiler("⊕", sortCompilerToAlmostSortCompiler(sorts.ListCompileSum)).
		WithListCompiler("⊗", sortCompilerToAlmostSortCompiler(sorts.ListCompileProd)).
		WithListCompiler("=>", el_sorts.ListCompileLambda).
		WithListCompiler("let", el_sorts.ListCompileLet).
		WithListCompiler("match", el_sorts.ListCompileMatch("exact")).
		finalize()
}

// finalize - just for reset syntax to be neat
func (ctx Context) finalize() Context {
	return ctx
}
