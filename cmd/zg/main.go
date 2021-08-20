package main

import (
	"os"

	"github.com/maruel/subcommands"
)

func main() {
	app := &subcommands.DefaultApplication{
		Name:  "zerogame",
		Title: "zerogame",
		Commands: []*subcommands.Command{
			subcommands.CmdHelp,
			CmdFeed(),
			CmdInstall(),
			CmdUninstall(),
			CmdRun(),
		},
	}

	os.Exit(subcommands.Run(app, os.Args[1:]))
}
