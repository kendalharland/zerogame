package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kendalharland/zerogame"
	"github.com/maruel/subcommands"
)

const (
	defaultFeedPath = "feed.json"
)

func CmdFeed() *subcommands.Command {
	return &subcommands.Command{
		UsageLine: "feed",
		ShortDesc: "generates a new feed",
		LongDesc:  "generates a new feed",
		CommandRun: func() subcommands.CommandRun {
			c := &cmdFeed{}
			c.Flags.StringVar(&c.feedPath, "o", defaultFeedPath, "where to write the output feed")
			return c
		},
	}
}

type cmdFeed struct {
	subcommands.CommandRunBase

	feedPath string
	feed     zerogame.Feed
}

func (c *cmdFeed) Run(a subcommands.Application, _ []string, _ subcommands.Env) int {
	if err := c.execute(context.Background()); err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

func (c *cmdFeed) execute(ctx context.Context) error {
	p := prompt{
		stdin:  os.Stdin,
		stderr: os.Stderr,
		stdout: os.Stdout,
	}

	if err := c.getFeedProperties(&p, &c.feed); err != nil {
		return err
	}

	if c.feedPath == "" {
		fp, err := p.ReadNonEmptyString("Enter the output feed filename: ")
		if err != nil {
			return err
		}
		c.feedPath = fp
	}

	feedBytes, err := json.MarshalIndent(c.feed, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.feedPath, feedBytes, 0755); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Feed was written to %s!\n", c.feedPath)
	return nil
}

func (c *cmdFeed) getFeedProperties(p *prompt, f *zerogame.Feed) error {
	name, err := p.ReadNonEmptyString("Enter the feed name: ")
	if err != nil {
		return err
	}
	f.Name = name

	version, err := p.ReadNonEmptyString("Enter the feed version: ")
	if err != nil {
		return err
	}
	f.Version = version

	archiveURL, err := p.ReadURL("Enter the feed archive URL: ")
	if err != nil {
		return err
	}
	f.ArchiveURL = archiveURL
	f.ArchiveType = "zip"

	gpgSignatureURL, err := p.ReadString("Enter the GPG signature URL (optional): ")
	if err != nil {
		return err
	}
	f.GPGSignatureURL = gpgSignatureURL
	return nil
}
