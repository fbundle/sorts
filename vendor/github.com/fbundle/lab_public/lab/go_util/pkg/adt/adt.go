package adt

type Prod2[T1 any, T2 any] struct {
	t1 T1
	t2 T2
}

func (p Prod2[T1, T2]) Unwrap() (t1 T1, t2 T2) {
	return p.t1, p.t2
}

func NewProd2[T1 any, T2 any](t1 T1, t2 T2) Prod2[T1, T2] {
	return Prod2[T1, T2]{
		t1: t1,
		t2: t2,
	}
}

type Sum2[T1 any, T2 any] struct {
	val any
}

func (s Sum2[T1, T2]) Unwrap1(val *T1) bool {
	if v, ok := s.val.(T1); ok {
		if val != nil {
			*val = v
		}
		return true
	} else {
		return false
	}
}
func (s Sum2[T1, T2]) Unwrap2(val *T2) bool {
	if v, ok := s.val.(T2); ok {
		if val != nil {
			*val = v
		}
		return true
	} else {
		return false
	}
}

func NewSum2[T1 any, T2 any](val any) Sum2[T1, T2] {
	okCount := 0
	if _, ok := val.(T1); ok {
		okCount++
	}
	if _, ok := val.(T2); ok {
		okCount++
	}
	if okCount == 0 {
		panic("type_error")
	}
	return Sum2[T1, T2]{val: val}
}
