package main

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/util"
	"github.com/urfave/cli"
)

const (
	blocktypeContent = "Content: "
)

func BlocksCmd(ctx *cli.Context) error {
	fixity, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	b, err := fixity.Head()
	if err != nil {
		return err
	}

	longHashes := ctx.Bool("long-hashes")

	blockPad := util.NewPad(" ", 4)
	printBlock(blockPad, longHashes, b)

	for i := 0; i < ctx.Int("limit") && b.PreviousBlockHash != ""; i++ {
		p, err := b.PreviousBlock()
		if err != nil {
			return err
		}
		b = p

		printBlock(blockPad, longHashes, b)
	}

	return nil
}

func printBlock(blockPad *util.AdjustPad, longHashes bool, b fixity.Block) {
	var blockHash string
	if longHashes {
		blockHash = b.BlockHash
	} else {
		blockHash = b.BlockHash[len(b.BlockHash)-8:]
	}

	fmt.Printf("%s %s\n",
		color.GreenString(blockPad.LeftPad(strconv.Itoa(b.Block))),
		color.YellowString(blockHash),
	)
}
