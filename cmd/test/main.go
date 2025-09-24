package main

import "fmt"

type Code interface {
	Haha()
}
type Sort interface {
	Hehe()
}

type Box[T any] struct {
	value T
}

func (b Box[T]) TypeName() string {
	switch zero[T]().(type) {
	case Code:
		return "code"
	case Sort:
		return "sort"
	default:
		return "unk"
	}
}

func zero[T any]() any {
	var z T
	return z
}

func main() {
	b1 := Box[Code]{}
	b2 := Box[Sort]{}
	fmt.Println(b1.TypeName(), b2.TypeName())
}
