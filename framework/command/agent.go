package command

import (
	"context"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/urfave/cli/v2"
)

type Agent struct {
	wrappers []*Wrapper
}

func NewAgent(wrappers ...*Wrapper) *Agent {
	return &Agent{
		wrappers: wrappers,
	}
}

func NewAgentOption(agent *Agent) base.Option {
	return base.NewOption(agent)
}

func (agent *Agent) AppendWrappers(wrappers ...*Wrapper) {
	agent.wrappers = append(agent.wrappers, wrappers...)
}

func (agent *Agent) Run(ctx context.Context, arguments []string) error {
	app := &cli.App{
		Commands: make([]*cli.Command, 0),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  cliFlag,
				Value: true,
			},
			&cli.BoolFlag{
				Name:  debugFlag,
				Value: true,
			},
		},
	}

	for _, wrapper := range agent.wrappers {
		app.Commands = append(app.Commands, wrapper.cliCommand())
	}

	return app.RunContext(ctx, arguments)
}
