package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func DeleteCmd(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	if id == "" {
		return cli.ShowCommandHelp(ctx, "delete")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	if err := fixi.Delete(id); err != nil {
		return err
	}

	// TODO(leeola): what is useful information to display here?
	// Affected hashes? afftected blocks? etc
	// Careful not to imply the content has been immediately removed.
	fmt.Println("id %q has been removed from the blockchain", id)

	return nil
}
