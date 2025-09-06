package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
)

func Name(s WithName) string {
	return s.nameAttr()
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
	WithName
	sortAttr() sortAttr
}

type sortAttr struct {
	level     int                     // universe Level
	name      string                  // every WithSort is identified with a Name (string)
	parent    WithSort                // (or Type) every WithSort must have a Parent
	lessEqual func(dst WithSort) bool // a partial order on sorts (subtype)
}

type WithName interface {
	nameAttr() string
}
