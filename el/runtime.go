package el

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Exec func(frame Frame, expr form.Form) (Frame, sorts.Sort, form.Form, error)

type _value struct {
	Sort  sorts.Sort
	Value form.Form
}

type Frame struct {
	Dict ordered_map.OrderedMap[form.Term, any] // Term -> Union[_value, Exec]
}

func (frame Frame) GetExec(key form.Term) (Exec, bool) {
	if o, ok := frame.Dict.Get(key); ok {
		if exec, ok := o.(Exec); ok {
			return exec, true
		}
	}
	return nil, false
}

func (frame Frame) SetExec(key form.Term, exec Exec) Frame {
	return Frame{
		Dict: frame.Dict.Set(key, exec),
	}
}

func (frame Frame) GetValue(key form.Term) (sorts.Sort, form.Form, bool) {
	if o, ok1 := frame.Dict.Get(key); ok1 {
		if val, ok2 := o.(_value); ok2 {
			return val.Sort, val.Value, true
		}
	}

	keyStr := string(key)

	// universes
	if strings.HasPrefix(keyStr, "U_") {
		levelStr := strings.TrimPrefix(keyStr, "U_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, nil, false
		}
		return sorts.NewAtom(level+1, keyStr, nil), key, true
	}
	return nil, nil, false
}

func (frame Frame) SetValue(key form.Term, sort sorts.Sort, form form.Form) Frame {
	return Frame{
		Dict: frame.Dict.Set(key, _value{Sort: sort, Value: form}),
	}
}

func (frame Frame) Eval(expr form.Form) (Frame, sorts.Sort, form.Form, error) {
	switch expr := expr.(type) {
	case form.Term:
		sort, value, ok := frame.GetValue(expr)
		if !ok {
			return frame, nil, nil, errors.New("variable not found: " + string(expr))
		}
		return frame, sort, value, nil
	case form.List:
		if len(expr) == 0 {
			return frame, nil, nil, errors.New("empty list")
		}
		if cmd, ok := expr[0].(form.Term); ok {
			if exec, ok := frame.GetExec(cmd); ok {
				// built-in function
				return exec(frame, expr)
			}
		}
		// normal function call
		if len(expr) != 2 {
			return frame, nil, nil, errors.New("regular function call must have exactly 1 argument")
		}
		cmdExpr, argExpr := expr[0], expr[1]

		frame, argSort, argValue, err := frame.Eval(argExpr)
		if err != nil {
			return frame, nil, nil, err
		}

		frame, cmdSort, cmdValue, err := frame.Eval(cmdExpr)
		if err != nil {
			return frame, nil, nil, err
		}

		if cmdValue != cmdExpr {
			return frame.Eval(form.List{cmdValue, argValue})
		}

		// cmd is undef - symbolically evaluate
		// type-check
		cmdSortParent := sorts.Parent(cmdSort)
		parentArrow, ok := cmdSortParent.(sorts.Arrow)
		if !ok {
			return frame, nil, nil, fmt.Errorf("expected arrow: %T", cmdSortParent)
		}
		if !sorts.TermOf(argSort, parentArrow.A) {
			return frame, nil, nil, fmt.Errorf("expected arg of type %v, got %v", parentArrow.A, argSort)
		}
		output := form.List{cmdValue, argValue}
		outputType := parentArrow.B
		return frame, sorts.NewAtom(sorts.Level(outputType)-1, form.String(output), outputType), output, nil
	default:
		return frame, nil, nil, errors.New("unknown expression: " + fmt.Sprintf("%T", expr))
	}
}
