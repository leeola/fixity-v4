package main

import (
	"io"
	"os"

	"github.com/urfave/cli"
)

func blobCommand(c *cli.Context) error {
	h := c.Args().Get(0)

	if h == "" {
		return cli.ShowCommandHelp(c, "blob")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	rc, err := client.GetBlob(h)
	if err != nil {
		return err
	}
	defer rc.Close()

	io.Copy(os.Stdout, rc)

	return nil
}
