package sorts

import "github.com/fbundle/sorts/adt"

func MustAtom(level int, name string, parent Sort) Sort {
	var sort Sort
	if ok := Atom(level, name, parent).Unwrap(&sort); !ok {
		panic("type_error")
	}
	return sort
}

func Atom(level int, name string, parent Sort) adt.Option[Sort] {
	if parent != nil && level+1 != parent.Level() {
		// if parent is specified, then its level must be valid
		return adt.None[Sort]()
	}
	return adt.Some[Sort](atom{
		level:  level,
		name:   name,
		parent: parent,
	})
}

// atom - representing all primitive sorts
// level 1: Int, Bool
// level 2: TypeName
type atom struct {
	level  int
	name   string
	parent Sort
}

func (s atom) Level() int {
	return s.level
}

func (s atom) String() string {
	return s.name
}

func (s atom) Parent() Sort {
	if s.parent != nil {
		return s.parent
	}
	// default parent - must have this to avoid infinity
	return atom{
		level: s.level + 1,
		name:  DefaultSortName,
	}
}

func (s atom) Length() int {
	return 1
}

func (s atom) LessEqual(dst Sort) bool {
	if s.Length() != dst.Length() || s.Level() != dst.Level() {
		return false
	}
	return le(s.String(), dst.String())
}

func (s atom) prepend(param Sort) Sort {
	return arrow{
		params: adt.MustNonEmpty([]Sort{param}),
		body:   s,
	}
}
