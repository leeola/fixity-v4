package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload"
	"github.com/leeola/fixity/util/clijson"
	"github.com/urfave/cli"
)

func WriteCmd(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		return cli.ShowCommandHelp(ctx, "write")
	}

	filePath := ctx.String("file")
	if filePath != "" {
		return errors.New("--file not implemented yet")
	}

	blobStr := ctx.String("blob")
	if blobStr != "" {
		return errors.New("--blob not implemented yet")
	}

	fixi, err := autoload.LoadFixity(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	b, err := clijson.CliJson(ctx.Args())
	if err != nil {
		return err
	}

	c := fixity.Commit{}
	j := fixity.Json{
		Json: json.RawMessage(b),
	}

	hashes, err := fixi.Write(c, j, nil)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(hashes, "\n"))

	return nil
}
