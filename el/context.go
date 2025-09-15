package el

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
	"github.com/fbundle/sorts/universe"
)

var TypeErr = fmt.Errorf("type_error")

var _ sorts.Context = Context{}

type Context struct {
	frame               ordered_map.OrderedMap[form.Name, sorts.Sort]
	universe            universe.SortUniverse
	listCompiler        ordered_map.OrderedMap[form.Name, sorts.ListCompileFunc]
	defaultListCompiler sorts.ListCompileFunc
	mode                sorts.Mode
}

func (ctx Context) WithMode(mode sorts.Mode) Context {
	ctx.mode = mode
	return ctx
}

func (ctx Context) Reset() Context {
	return Context{
		frame: ordered_map.OrderedMap[form.Name, sorts.Sort]{},
		universe: universe.SortUniverse{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
		},
		listCompiler:        ordered_map.OrderedMap[form.Name, sorts.ListCompileFunc]{},
		defaultListCompiler: sorts.ListCompileBeta,
		mode:                sorts.ModeComp,
	}.
		WithListCompiler("->", sorts.ListCompileArrow).
		WithListCompiler("⊕", sorts.ListCompileSum).
		WithListCompiler("⊗", sorts.ListCompileProd).
		WithListCompiler("=>", sorts.ListCompileLambda(":")).
		WithListCompiler("inh", sorts.ListCompileInhabitant).
		WithListCompiler("let", sorts.ListCompileLet(":=")).
		WithListCompiler("match", sorts.ListCompileMatch("=>", "_")).
		WithListCompiler("inspect", sorts.ListCompileInspect).
		WithListCompiler("type", sorts.ListCompileType).
		finalize()
}

// finalize - just for reset syntax to be neat
func (ctx Context) finalize() Context {
	return ctx
}
