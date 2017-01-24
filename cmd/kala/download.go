package main

import (
	"io"
	"os"

	"github.com/leeola/errors"
	"github.com/urfave/cli"
)

func downloadCommand(c *cli.Context) error {
	var w io.Writer
	if c.Bool("stdout") {
		w = os.Stdout
		// Continue here, add upload --stdin junk so it uploads files by default
	} else if p := c.Args().Get(1); p != "" {
		f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	if w == nil {
		cli.ShowSubcommandHelp(c)
		return errors.New("error: missing path")
	}

	h := c.Args().Get(0)
	if h == "" {
		cli.ShowSubcommandHelp(c)
		return errors.New("error: missing hash")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	rc, err := client.Download(h)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(w, rc); err != nil {
		return err
	}

	return nil
}
