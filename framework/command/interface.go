package command

import "github.com/urfave/cli/v2"

type Interface interface {
	Usage() string
	Command() string
	Arguments() []cli.Flag
	Run(*cli.Context) error
}
