package main

import (
	"io"
	"os"

	"github.com/urfave/cli"
)

func blobCommand(c *cli.Context) error {
	if c.Bool("upload") {
		return uploadBlobCommand(c)
	}
	return getBlobCommand(c)
}

func getBlobCommand(c *cli.Context) error {
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

func uploadBlobCommand(c *cli.Context) error {
	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	h, err := client.PostBlob(os.Stdin)
	if err != nil {
		return err
	}

	Printlnf(h)

	return nil
}
