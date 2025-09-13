package el2

import (
	"github.com/fbundle/sorts/sorts"
	"github.com/fbundle/sorts/universe"
)

func newUniverse() universe.Universe {
	u, err := universe.New("U", "Unit", "Any")
	if err != nil {
		panic(err)
	}

	err = u.NewParseListRule("->", sorts.ParseListArrow("->"))
	if err != nil {
		panic(err)
	}
	return u

}
