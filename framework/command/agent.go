package command

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/urfave/cli/v2"
	"os"
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

func (agent *Agent) OsArgsContains(s string) bool {
	for _, arg := range os.Args {
		if arg == s {
			return true
		}
	}
	return false
}

func (agent *Agent) MustRun() bool {
	return agent.OsArgsContains(fmt.Sprintf("--%s", cliFlag))
}

func (agent *Agent) MustDebug() bool {
	return agent.OsArgsContains(fmt.Sprintf("--%s", debugFlag))
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
				Value: false,
			},
		},
	}

	for _, wrapper := range agent.wrappers {
		app.Commands = append(app.Commands, wrapper.cliCommand())
	}

	return app.RunContext(ctx, arguments)
}
