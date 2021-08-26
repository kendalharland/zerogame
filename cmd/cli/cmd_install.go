package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/kendalharland/zerogame"
	"github.com/maruel/subcommands"
)

func CmdInstall() *subcommands.Command {
	return &subcommands.Command{
		UsageLine: "install feed_url",
		ShortDesc: "installs an archive from a feed URL",
		LongDesc:  "installs an archive from a feed URL",
		CommandRun: func() subcommands.CommandRun {
			c := &cmdInstall{}
			c.Flags.BoolVar(&c.noVerify, "noverify", false, "Disables Feed verification")
			c.Flags.BoolVar(&c.noCache, "nocache", false, "Forces downloading the feed even if it exists locally")
			return c
		},
	}
}

type cmdInstall struct {
	subcommands.CommandRunBase

	noVerify bool
	noCache  bool
}

func (c *cmdInstall) Run(a subcommands.Application, _ []string, _ subcommands.Env) int {
	if err := c.execute(context.Background()); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func (c *cmdInstall) execute(ctx context.Context) error {
	if c.Flags.NArg() != 1 {
		return errors.New("expected one argument")
	}

	home, _ := os.UserHomeDir()
	workspace := filepath.Join(home, ".config", "zerogame")
	db := zerogame.NewFeedRecordDB(filepath.Join(workspace, "db"))
	repo := zerogame.NewFeedRepository(workspace, db)

	return zerogame.InstallFeed(ctx, c.Flags.Arg(0), repo)
}
