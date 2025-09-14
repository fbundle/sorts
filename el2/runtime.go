package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/frame"
	"github.com/fbundle/sorts/el2/parser"
	"github.com/fbundle/sorts/el2/sort_universe"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

type Runtime struct {
	frame        el2_frame.Frame
	sortUniverse el2_sort_universe.SortUniverse
}

func NewRuntime() Runtime {
	return Runtime{
		frame: el2_frame.Frame{},
		sortUniverse: el2_sort_universe.SortUniverse{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
		},
	}
}

func NewParser(r Runtime) el2_parser.Parser {
	p := el2_parser.Parser{
		ParseName: nil,
	}
	p.ParseName = func(name form.Name) sorts.Sort {
		if sort, ok := r.frame.Get(name); ok {
			return sort
		}
		if sort, ok := r.sortUniverse.ParseBuiltin(name); ok {
			return sort
		}
		panic(fmt.Errorf("name_not_found: %s", name))
	}

	return p.
		WithListParser("->", toAlmostSortListParser(sorts.ListParseArrow)).
		WithListParser("⊕", toAlmostSortListParser(sorts.ListParseSum)).
		WithListParser("⊗", toAlmostSortListParser(sorts.ListParseProd)).
		WithListParser("=>", el2_almost_sort.ListParseLambda).
		WithListParser("let", el2_almost_sort.ListParseLet("undef")).
		WithListParser("match", el2_almost_sort.ListParseMatch("exact"))
}
