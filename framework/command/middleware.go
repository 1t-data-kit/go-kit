package command

import (
	"github.com/urfave/cli/v2"
	"reflect"
)

type Middleware func(cli.ActionFunc) cli.ActionFunc

func applyMiddleware(f cli.ActionFunc, middlewares ...Middleware) cli.ActionFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		apply := middlewares[i]
		if apply == nil || reflect.ValueOf(apply).IsNil() {
			continue
		}
		f = apply(f)
	}
	return f
}
