package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) match(condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
	frame, compSort, compValue, err := frame.Resolve(comp)
	if err != nil {
		return frame, false, err
	}
	if compSort == condSort {
		if String(compValue) == String(condValue) {
			// TODO - improve match - currently we just have exact match
			return frame, true, nil
		}
	}
	return frame, false, nil
}
