package main

import (
	"os"

	"github.com/urfave/cli"
)

func uploadCommand(c *cli.Context) error {
	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	hashes, err := client.Upload("file", os.Stdin)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}
