package sorts_parser

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/slices_util"
	"github.com/fbundle/sorts/sorts"
)

func (p Parser) withListParseFunc(cmd Name, parseFunc ListParseFunc) Parser {
	p.listParseFuncMap = p.listParseFuncMap.Set(cmd, parseFunc)
	return p
}

func (p Parser) finalize() Parser {
	if p.BuiltinNameParseFunc == nil {
		p.BuiltinNameParseFunc = func(name sorts.Name) (sorts.Code, bool) {
			return nil, false
		}
	}
	return p
}

func (p Parser) Init() Parser {
	p.finalNameParseFunc = func(name sorts.Name) sorts.Code {
		return sorts.Symbol{Name: name}
	}
	p.finalListParseFunc = func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code {
		err := compileErr(list, []string{
			"cmd",
			"arg1",
			"...",
			"argN",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}
		output := parse(list[0])
		args := slices_util.Map(list[1:], func(form sorts.Form) sorts.Code {
			return parse(form)
		})
		slices_util.ForEach(args, func(arg sorts.Code) {
			output = sorts.Beta{
				Cmd: output,
				Arg: arg,
			}
		})
		return output
	}
	const (
		ArrowCmd sorts.Name = "->"
	)

	return p.
		withListParseFunc(sorts.TypeCmd, func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code {
			err := compileErr(list, []string{string(sorts.TypeCmd), "value"})
			if len(list) != 1 {
				panic(err)
			}
			return sorts.Type{
				Value: parse(list[0]),
			}
		}).
		withListParseFunc(sorts.InhabitCmd, func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code {
			err := compileErr(list, []string{string(sorts.InhabitCmd), "type"})
			if len(list) != 1 {
				panic(err)
			}
			typeCode := parse(list[0])
			return sorts.Inhabited{
				Type: typeCode,
			}
		}).
		withListParseFunc(sorts.PiCmd, func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code {
			err := compileErr(list, []string{
				string(sorts.PiCmd),
				makeForm(sorts.AnnotCmd, "name1", "type1"),
				"...",
				makeForm(sorts.AnnotCmd, "nameN", "typeN"),
				"body",
			}, "where N >= 0")
			if len(list) < 1 {
				panic(err)
			}
			params := slices_util.Map(list[:len(list)-1], func(form sorts.Form) sorts.Annot {
				return compileAnnot(parse, mustType[sorts.List](err, form)[1:])
			})
			output := parse(list[len(list)-1])
			slices_util.ForEach(slices_util.Reverse(params), func(param sorts.Annot) {
				output = sorts.Pi{
					Param: param,
					Body:  output,
				}
			})
			return output
		}).
		withListParseFunc(ArrowCmd, func(parse func(form sorts.Form) sorts.Code, list sorts.List) sorts.Code {
			// make builtin like succ
			// e.g. if arrow is Nat -> Nat
			// then its lambda is
			// (x: Nat) => Nat
			// or some mechanism to introduce arrow type from pi type
			// TODO - probably we don't need this anymore
			panic("not implemented")
		}).
		withListParseFunc(sorts.LetCmd, func(parse func(form form.Form) sorts.Code, list form.List) Code {
			err := compileErr(list, []string{
				string(sorts.LetCmd),
				makeForm(sorts.BindingCmd, "name1", "value1"),
				"...",
				makeForm(sorts.BindingCmd, "nameN", "valueN"),
				"body",
			}, "where N >= 0")
			if len(list) < 1 {
				panic(err)
			}
			bindings := slices_util.Map(list[:len(list)-1], func(form Form) sorts.Binding {
				return compileBinding(parse, mustType[List](err, form)[1:])
			})
			output := parse(list[len(list)-1])
			slices_util.ForEach(slices_util.Reverse(bindings), func(binding sorts.Binding) {
				output = sorts.Let{
					Binding: binding,
					Body:    output,
				}
			})
			return output
		}).
		finalize()
}
