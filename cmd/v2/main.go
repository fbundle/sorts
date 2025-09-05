package main

import (
	"fmt"

	"github.com/fbundle/sorts/sorts/v2"
)

const (
	defaultName  = "type"
	unitTypeName = "unit"
	anyTypeName  = "any"

	dataLevel = 0
	typeLevel = 1
)

func chain(ss sorts.SortSystem, sortList ...sorts.Sort) sorts.Sort {
	if len(sortList) == 0 {
		return nil
	}
	final := sortList[len(sortList)-1]
	for i := len(sortList) - 2; i >= 0; i-- {
		param := sortList[i]
		final = ss.Arrow(param, final)
	}
	return final
}

func weakSort(ss sorts.SortSystem, level int, numParams int) sorts.Sort {
	if numParams < 0 {
		panic("type_error")
	}
	anyType := ss.Atom(level, anyTypeName)
	unitType := ss.Atom(level, unitTypeName)

	sortList := make([]sorts.Sort, 0, numParams+1)
	for i := 0; i < numParams; i++ {
		sortList = append(sortList, anyType)
	}
	sortList = append(sortList, unitType)
	return chain(ss, sortList...)
}

func strongSort(ss sorts.SortSystem, level int, numParams int) sorts.Sort {
	if numParams < 0 {
		panic("type_error")
	}
	anyType := ss.Atom(level, anyTypeName)
	unitType := ss.Atom(level, unitTypeName)

	sortList := make([]sorts.Sort, 0, numParams+1)
	for i := 0; i < numParams; i++ {
		sortList = append(sortList, unitType)
	}
	sortList = append(sortList, anyType)
	return chain(ss, sortList...)
}

func main() {
	fmt.Printf("anything can be cast into [%s]\n", anyTypeName)
	fmt.Printf("[%s] can be cast into anything\n", unitTypeName)

	ss := sorts.NewSortSystem(defaultName, sorts.WithInitialName(unitTypeName), sorts.WithTerminalName(anyTypeName))

	intType := ss.Atom(typeLevel, "int")
	boolType := ss.Atom(typeLevel, "bool")
	stringType := ss.Atom(typeLevel, "string")

	intBoolType := chain(ss, intType, boolType)
	intIntIntType := chain(ss, intType, intType, intType)

	weak1 := weakSort(ss, typeLevel, 1)
	strong1 := strongSort(ss, typeLevel, 1)
	weak3 := weakSort(ss, typeLevel, 3)
	strong3 := strongSort(ss, typeLevel, 3)

	ss = ss.AddRule("bool", "int") // cast bool -> int
	fmt.Println("[bool] can be cast into [int]")

	printSorts(stringType, intIntType, intIntIntType, weak1, strong1, weak3, strong3)

	printCast(boolType, intType)
	printCast(intType, boolType)
	printCast(weak3, intIntType)
	printCast(strong3, intIntType)
}
