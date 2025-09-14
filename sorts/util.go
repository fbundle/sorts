package sorts

import "github.com/fbundle/sorts/form"

func LeastUpperBound(a SortAttr, H form.Name, sorts ...Sort) Sort {
	removeAny := func() int {
		for i := 0; i < len(sorts); i++ {
			for j := 0; j < len(sorts); j++ {
				if i == j {
					continue
				}
				if a.LessEqual(sorts[i], sorts[j]) {
					return i
				}
			}
		}
		return -1
	}
	for {
		i := removeAny()
		if i < 0 {
			break
		}
		sorts[i] = sorts[len(sorts)-1]
		sorts = sorts[:len(sorts)-1]
	}

	output := sorts[0]
	for i := 1; i < len(sorts); i++ {
		output = Sum{
			H: H,
			A: sorts[i],
			B: output,
		}
	}
	return output
}

func must(a SortAttr) mustSortAttr {
	return mustSortAttr{a}
}

type mustSortAttr struct {
	a SortAttr
}

func (m mustSortAttr) lessEqual(x Sort, y Sort) {
	if !m.a.LessEqual(x, y) {
		panic(TypeErr)
	}
}

func (m mustSortAttr) termOf(x Sort, X Sort) {
	if !m.a.LessEqual(m.a.Parent(x), X) {
		panic(TypeErr)
	}
}
