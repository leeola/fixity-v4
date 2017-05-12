package main

import (
	"errors"

	"github.com/urfave/cli"
)

func SearchCmd(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		return cli.ShowCommandHelp(ctx, "write")
	}

	return errors.New("not implemented")
}
