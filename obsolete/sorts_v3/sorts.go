package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial almost_sort
	TerminalName = "any"  // terminal almost_sort
)

func Name(s any) string {
	if s == nil {
		return "nil"
	}
	switch s := s.(type) {
	case Sort:
		return s.sortAttr().name
	case Dependent:
		return s.Name
	default:
		panic("type_error")
	}
}

func Level(s Sort) int {
	return s.sortAttr().level
}

func Parent(s Sort) Sort {
	return s.sortAttr().parent
}

func SubTypeOf(x Sort, y Sort) bool {
	if Level(x) == Level(y) && (Name(x) == InitialName || Name(y) == TerminalName) {
		return true
	}

	return x.sortAttr().lessEqual(y)
}
func TermOf(x Sort, X Sort) bool {
	X1 := Parent(x)
	return SubTypeOf(X1, X)
}

type Sort interface {
	sortAttr() sortAttr
}

type sortAttr struct {
	name      string              // every Sort is identified with a Name (string)
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}
