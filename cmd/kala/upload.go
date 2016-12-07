package main

import (
	"io"
	"os"

	"github.com/leeola/errors"
	"github.com/urfave/cli"
)

func uploadCommand(c *cli.Context) error {
	var r io.Reader
	if c.Bool("stdin") {
		r = os.Stdin
		// Continue here, add upload --stdin junk so it uploads files by default
	} else if p := c.Args().First(); p != "" {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	if r == nil {
		cli.ShowSubcommandHelp(c)
		return errors.New("error: missing path")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	hashes, err := client.Upload("file", r)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}
