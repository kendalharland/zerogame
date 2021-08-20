package main

import (
	"context"
	"errors"
	"log"

	"github.com/maruel/subcommands"
)

func CmdRun() *subcommands.Command {
	return &subcommands.Command{
		UsageLine:  "run",
		ShortDesc:  "runs an archive from a feed URL",
		LongDesc:   "runs an archive from a feed URL",
		CommandRun: func() subcommands.CommandRun { return &cmdRun{} },
	}
}

type cmdRun struct {
	subcommands.CommandRunBase
}

func (c *cmdRun) Run(a subcommands.Application, _ []string, _ subcommands.Env) int {
	if err := c.execute(context.Background()); err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

func (c *cmdRun) execute(ctx context.Context) error {
	return errors.New("unimplemented")
}
