package el2

import (
	"github.com/fbundle/sorts/el2/el_almost_sort"
	el_frame "github.com/fbundle/sorts/el2/frame"
	el_sort_universe "github.com/fbundle/sorts/el2/sort_universe"
)

type Runtime struct {
	frame        el_frame.Frame
	sortUniverse el_sort_universe.SortUniverse
	parser       Parser
}

func newRuntime() Runtime {
	frame := el_frame.Frame{}
	sortUniverse := el_sort_universe.SortUniverse{
		InitialTypeName:  "Unit",
		TerminalTypeName: "Any",
	}
	parser := Parser{}

	if sort, ok := r.frame.Get(node); ok {
		return el_almost_sort.ActualSort{sort}
	}
	if sort, ok := r.sortUniverse.ParseBuiltin(node); ok {
		return el_almost_sort.ActualSort{sort}
	}
	panic("name_not_found")
}
