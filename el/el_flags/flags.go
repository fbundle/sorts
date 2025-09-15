package el_flags

import "os"

type Mode string

const (
	ModeComp  Mode = "COMP"  // type checking
	ModeEval  Mode = "EVAL"  // type checking and evaluation
	ModeDebug Mode = "DEBUG" // type checking and print everything
)

func GetMode() Mode {
	modeEnv := Mode(os.Getenv("EL_MODE"))

	modes := map[Mode]struct{}{
		ModeComp:  {},
		ModeEval:  {},
		ModeDebug: {},
	}
	if _, ok := modes[modeEnv]; ok {
		return modeEnv
	}
	return ModeComp
}
