package main

import (
	"fmt"

	"github.com/fbundle/sorts/sorts/sorts_v2"
)

const (
	defaultName  = "type"
	unitTypeName = "unit"
	anyTypeName  = "any"

	dataLevel = 0
	typeLevel = 1
)

func printCast(type1 sorts.Sort, type2 sorts.Sort) {
	ok := type1.LessEqual(type2)

	if ok {
		fmt.Printf("type [%s] CAN be cast into [%s]\n", type1.Name(), type2.Name())
	} else {
		fmt.Printf("type [%s] CANNOT be cast into [%s]\n", type1.Name(), type2.Name())
	}
}
func printSorts(sortList ...sorts.Sort) {
	for _, sort := range sortList {
		fmt.Printf("[%s] is of type [%s]\n", sort.Name(), sort.Parent().Sort().Name())
	}
}
func chain(ss sorts.SortSystem, sortList ...sorts.Sort) sorts.Sort {
	if len(sortList) == 0 {
		return nil
	}
	final := sortList[len(sortList)-1]
	for i := len(sortList) - 2; i >= 0; i-- {
		arg := sortList[i]
		if ok := ss.Arrow(arg, final).Unwrap(&final); !ok {
			panic("type_error")
		}
	}
	return final
}

func mustAtom(ss sorts.SortSystem, level int, name string, parent sorts.Sort) sorts.Sort {
	var atom sorts.Sort
	if ok := ss.Atom(level, name, parent).Unwrap(&atom); !ok {
		panic("type_error")
	}
	return atom
}

func weakSort(ss sorts.SortSystem, level int, numParams int) sorts.Sort {
	if numParams < 0 {
		panic("type_error")
	}
	anyType := mustAtom(ss, level, anyTypeName, nil)
	unitType := mustAtom(ss, level, unitTypeName, nil)

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
	anyType := mustAtom(ss, level, anyTypeName, nil)
	unitType := mustAtom(ss, level, unitTypeName, nil)

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

	intType := mustAtom(ss, typeLevel, "int", nil)
	boolType := mustAtom(ss, typeLevel, "bool", nil)
	stringType := mustAtom(ss, typeLevel, "string", nil)

	intBoolType := chain(ss, intType, boolType)
	intIntIntType := chain(ss, intType, intType, intType)

	weak1 := weakSort(ss, typeLevel, 1)
	strong1 := strongSort(ss, typeLevel, 1)
	weak2 := weakSort(ss, typeLevel, 2)
	strong2 := strongSort(ss, typeLevel, 2)

	ss = ss.AddRule("bool", "int") // cast bool -> int
	fmt.Println("[bool] can be cast into [int]")

	printSorts(stringType, intBoolType, intIntIntType, weak1, strong1, weak2, strong2)

	printCast(boolType, intType)
	printCast(intType, boolType)
	printCast(weak2, intIntIntType)
	printCast(strong2, intIntIntType)
}
