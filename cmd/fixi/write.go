package main

import (
	"errors"

	"github.com/urfave/cli"
)

func WriteCmd(ctx *cli.Context) error {
	filePath := ctx.String("file")
	if filePath != "" {
		return errors.New("--file not implemented yet")
	}

	blobStr := ctx.String("blob")
	if blobStr != "" {
		return errors.New("--blob not implemented yet")
	}

	return errors.New("not implemented")
}
