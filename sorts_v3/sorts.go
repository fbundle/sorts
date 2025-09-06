package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
)

func Name(s any) string {
	switch s := s.(type) {
	case WithSort:
		return s.sortAttr().name
	case Dependent:
		return s.Name
	default:
		panic("type_error")
	}
}

func Level(s WithSort) int {
	return s.sortAttr().level
}

func Parent(s WithSort) WithSort {
	return s.sortAttr().parent
}

func LessEqual(x WithSort, y WithSort) bool {
	return x.sortAttr().lessEqual(y)
}

type WithSort interface {
	sortAttr() sortAttr
}

type sortAttr struct {
	name      string                  // every WithSort is identified with a Name (string)
	level     int                     // universe Level
	parent    WithSort                // (or Type) every WithSort must have a Parent
	lessEqual func(dst WithSort) bool // a partial order on sorts (subtype)
}
