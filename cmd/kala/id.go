package main

import "github.com/urfave/cli"

func idCommand(c *cli.Context) error {
	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	id, err := client.NodeId()
	if err != nil {
		return err
	}

	Printlnf(id)

	return nil
}
