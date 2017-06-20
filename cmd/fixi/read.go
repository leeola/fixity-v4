package main

import (
	"errors"
	"io"
	"os"

	"github.com/leeola/fixity"
	"github.com/urfave/cli"
)

func ReadCmd(ctx *cli.Context) error {
	idOrHash := ctx.Args().Get(0)
	path := ctx.Args().Get(1)

	useStdout := ctx.Bool("stdout")
	if idOrHash == "" || (path == "" && !useStdout) {
		return cli.ShowCommandHelp(ctx, "read")
	}

	if path != "" && useStdout {
		return errors.New("error: path cannot be used with --stdout flag")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	var content fixity.Content
	if !ctx.Bool("hash") {
		content, err = fixi.Read(idOrHash)
	} else {
		content, err = fixi.ReadHash(idOrHash)
	}
	if err != nil {
		return err
	}

	rc, err := content.Read()
	if err != nil {
		return err
	}
	defer rc.Close()

	var out io.Writer
	if useStdout {
		out = os.Stdout
	} else {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	_, err = io.Copy(out, rc)
	return err
}
