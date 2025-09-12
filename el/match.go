package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) match(condSort sorts.Sort, cond Expr, comp Expr) (Frame, bool) {
	if comp, ok := comp.(Term); ok {
		// match all
		if _, _, ok := frame.Get(comp); ok {
			// exact match
			return frame, cond == comp
		} else {
			// normal match
			return frame.Set(comp, condSort, cond), true
		}
	}
	if comp, ok := comp.(FunctionCall); ok {
		if cond, ok := cond.(FunctionCall); ok {
			frame, cmdSort, _, err := frame.Resolve(comp.Cmd)
			if err != nil {
				return frame, false
			}
			frame, matched := frame.match(cmdSort, cond.Cmd, comp.Cmd)
			if !matched {
				return frame, false
			}
			frame, argSort, _, err := frame.Resolve(comp.Arg)
			if err != nil {
				return frame, false
			}
			return frame.match(argSort, cond.Arg, comp.Arg)
		}
	}
	return frame, false
}
