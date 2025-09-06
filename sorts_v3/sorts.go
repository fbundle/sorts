package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
)

func Name(s any) string {
	switch s := s.(type) {
	case WithSort:
		return s.attr().name
	case WithName:
		return s.attr().name
	default:
		panic("type_error")
	}
}

func Level(s WithSort) int {
	return s.attr().level
}

func Parent(s WithSort) WithSort {
	return s.attr().parent
}

func LessEqual(x WithSort, y WithSort) bool {
	return x.attr().lessEqual(y)
}

type WithSort interface {
	attr() sortAttr
}

type sortAttr struct {
	level     int                     // universe Level
	name      string                  // every WithSort is identified with a Name (string)
	parent    WithSort                // (or Type) every WithSort must have a Parent
	lessEqual func(dst WithSort) bool // a partial order on sorts (subtype)
}

type WithName interface {
	attr() nameAttr
}

type nameAttr struct {
	name string
}
