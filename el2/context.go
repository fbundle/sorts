package el2

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/el2/sort_universe"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
)

var TypeErr = fmt.Errorf("type_error")

var _ almost_sort_extra.Context = Context{}

type Context struct {
	frame        ordered_map.OrderedMap[form.Name, almost_sort.ActualSort]
	sortUniverse el2_sort_universe.SortUniverse
	listParsers  ordered_map.OrderedMap[form.Name, almost_sort_extra.ListCompileFunc]
}
