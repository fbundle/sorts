package sorts

import "github.com/fbundle/sorts/form"

func LeastUpperBound(a SortAttr, H form.Name, sorts ...Sort) Sort {
	ss := make(map[int]Sort)
	for i, s := range sorts {
		ss[i] = s
	}

	removeAny := func(ss map[int]Sort) int {
		for i1, s1 := range ss {
			for i2, s2 := range ss {
				if i1 == i2 {
					continue
				}
				if a.LessEqual(s1, s2) {
					return i1
				}
			}
		}
		return -1
	}
	for {
		i := removeAny(ss)
		if i < 0 {
			break
		}
		delete(ss, i)
	}

	sorts = make([]Sort, 0, len(ss))
	for _, s := range ss {
		sorts = append(sorts, s)
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
