package main

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/leeola/fixity/util/dyntabwriter"
	"github.com/urfave/cli"
)

func InfoCmd(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	if id == "" {
		return cli.ShowCommandHelp(ctx, "info")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	w := dyntabwriter.New(os.Stdout)
	defer w.Flush()
	w.Header("     ", "ID", "HASH", "SIZE", "AVG CHUNK")

	c, err := fixi.Read(id)
	if err != nil {
		return err
	}

	showHashes := ctx.Bool("full-hashes")

	blob, err := c.Blob()
	if err != nil {
		return err
	}

	w.Println("     ",
		color.GreenString(id),
		color.GreenString(sumHash(c.Hash, showHashes)),
		color.YellowString(strconv.Itoa(int(blob.Size))),
		color.YellowString(strconv.Itoa(int(blob.AverageChunkSize))),
	)

	for c.PreviousContentHash != "" {
		c, err = c.PreviousContent()
		if err != nil {
			return err
		}

		blob, err := c.Blob()
		if err != nil {
			return err
		}

		w.Println("     ",
			color.GreenString(id),
			color.GreenString(sumHash(c.Hash, showHashes)),
			color.YellowString(strconv.Itoa(int(blob.Size))),
			color.YellowString(strconv.Itoa(int(blob.AverageChunkSize))),
		)
	}

	// TODO(leeola): show summarized values of the total, deduped
	// storage use of the content.

	return nil
}
