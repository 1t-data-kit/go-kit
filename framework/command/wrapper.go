package command

import (
	"github.com/1t-data-kit/go-kit/base"
	"github.com/urfave/cli/v2"
)

type Wrapper struct {
	object      Object
	middlewares []Middleware
}

func NewWrapper(object Object, middlewares ...Middleware) *Wrapper {
	return &Wrapper{
		object:      object,
		middlewares: middlewares,
	}
}

func NewWrapperOption(wrapper *Wrapper) base.Option {
	return base.NewOption(wrapper)
}

func (wrapper *Wrapper) cliCommand() *cli.Command {
	return &cli.Command{
		Usage:  wrapper.object.Usage(),
		Name:   wrapper.object.Command(),
		Flags:  wrapper.object.Arguments(),
		Action: applyMiddleware(wrapper.object.Run, wrapper.middlewares...),
	}
}
