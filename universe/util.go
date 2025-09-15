package universe

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
