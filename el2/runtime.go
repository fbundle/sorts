package el2

import (
	"fmt"

	el_almost_sort "github.com/fbundle/sorts/el2/almost_sort"
	el_frame "github.com/fbundle/sorts/el2/frame"
	el_parser "github.com/fbundle/sorts/el2/parser"
	el_sort_universe "github.com/fbundle/sorts/el2/sort_universe"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

type Runtime struct {
	frame        el_frame.Frame
	sortUniverse el_sort_universe.SortUniverse
	parser       el_parser.Parser
}

func newRuntime() Runtime {
	r := Runtime{
		frame: el_frame.Frame{},
		sortUniverse: el_sort_universe.SortUniverse{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
		},
		parser: el_parser.Parser{
			ParseName: nil,
		},
	}
	r.parser.ParseName = func(name form.Name) sorts.Sort {
		if sort, ok := r.frame.Get(name); ok {
			return sort
		}
		if sort, ok := r.sortUniverse.ParseBuiltin(name); ok {
			return sort
		}
		panic(fmt.Errorf("name_not_found: %s", name))
	}
	return r
}

func DefaultRuntime() Runtime {
	r := newRuntime()
	r.parser = r.parser.
		NewListParser("->", toListParser(sorts.ListParseArrow("->"))).
		NewListParser("⊕", toListParser(sorts.ListParseSum("⊕"))).
		NewListParser("⊗", toListParser(sorts.ListParseProd("⊗"))).
		NewListParser("=>", el_almost_sort.ListParseLambda)
	return r
}
