package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
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

func LessEqual(x Sort, y Sort) bool {
	return x.sortAttr().lessEqual(y)
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
