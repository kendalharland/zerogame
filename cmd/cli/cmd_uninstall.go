package main

import (
	"context"
	"errors"
	"log"

	"github.com/maruel/subcommands"
)

func CmdUninstall() *subcommands.Command {
	return &subcommands.Command{
		UsageLine:  "uninstall",
		ShortDesc:  "Uninstalls an archive",
		LongDesc:   "Uninstalls an archive",
		CommandRun: func() subcommands.CommandRun { return &cmdUninstall{} },
	}
}

type cmdUninstall struct {
	subcommands.CommandRunBase
}

func (c *cmdUninstall) Run(a subcommands.Application, _ []string, _ subcommands.Env) int {
	if err := c.execute(context.Background()); err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

func (c *cmdUninstall) execute(ctx context.Context) error {
	return errors.New("unimplemented")
}
