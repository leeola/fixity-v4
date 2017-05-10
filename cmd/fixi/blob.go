package main

import (
	"errors"

	"github.com/urfave/cli"
)

func BlobCmd(ctx *cli.Context) error {
	h := ctx.Args().Get(0)
	if h == "" {
		return cli.ShowCommandHelp(ctx, "blob")
	}

	return errors.New("not implemented")
}
