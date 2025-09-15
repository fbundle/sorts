package el

import "github.com/fbundle/sorts/sorts"

func (ctx Context) Mode() sorts.Mode {
	return ctx.mode
}
