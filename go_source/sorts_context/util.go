package sorts_context

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
)

func getBuiltinLevel(builtinName form.Name, name form.Name) (int, bool) {
	nameStr := string(name)
	prefix := string(builtinName) + "_"
	if strings.HasPrefix(nameStr, prefix) {
		levelStr := strings.TrimPrefix(nameStr, prefix)
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return 0, false
		}
		return level, true
	}
	return 0, false
}

func setBuiltinLevel(builtinName form.Name, level int) form.Name {
	prefix := string(builtinName) + "_"
	return form.Name(prefix + strconv.Itoa(level))
}

func equalForm(s Form, d Form) bool {
	sName, ok1 := s.(Name)
	dName, ok2 := d.(Name)
	if ok1 && ok2 {
		return sName == dName
	}
	sList, ok1 := s.(List)
	dList, ok2 := d.(List)
	if ok1 && ok2 {
		if len(sList) != len(dList) {
			return false
		}
		for i := range sList {
			sVal, dVal := sList[i], sList[i]
			if !equalForm(sVal, dVal) {
				return false
			}
		}
		return true
	}
	return false
}
