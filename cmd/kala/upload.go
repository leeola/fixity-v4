package main

import (
	"io"
	"os"

	"github.com/leeola/errors"
	"github.com/urfave/cli"
)

func uploadCommand(c *cli.Context) error {
	metaChanges := argsToMetaChanges(c.Args())

	var r io.Reader
	if c.Bool("stdin") {
		r = os.Stdin
		// Continue here, add upload --stdin junk so it uploads files by default
	} else if p := c.String("file"); p != "" {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f

		// If the user didn't specify metadata to add themselves,
		// automatically set the filename and constenttype for an easy ux.
		if _, ok := metaChanges["filename"]; !ok {
			metaChanges.Set("filename", p)
		}
		if _, ok := metaChanges.GetContentType(); !ok {
			metaChanges.SetContentType("file")
		}
	}

	if r == nil {
		cli.ShowSubcommandHelp(c)
		return errors.New("error: must specify either file or stdin")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	hashes, err := client.Upload(r, metaChanges)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}
