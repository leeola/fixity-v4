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

	blockPad := util.NewPad(" ", 4)

	b, err := fixity.Head()
	if err != nil {
		return err
	}

	c, err := b.Content()
	if err != nil {
		return err
	}

	showBlockHashes := ctx.Bool("block-hashes")
	showContentHashes := ctx.Bool("content-hashes")

	block := blockPad.LeftPad(strconv.Itoa(b.Block))
	bHash := sumHash(b.BlockHash, showBlockHashes)
	cHash := sumHash(b.ContentHash, showContentHashes)

	printBlock(block, bHash, cHash, c.Id, c.IndexedFields)

	for i := 0; i < ctx.Int("limit") && b.PreviousBlockHash != ""; i++ {
		b, err = b.PreviousBlock()
		if err != nil {
			return err
		}

		c, err = b.Content()
		if err != nil {
			return err
		}

		block := blockPad.LeftPad(strconv.Itoa(b.Block))
		bHash := sumHash(b.BlockHash, showBlockHashes)
		cHash := sumHash(b.ContentHash, showContentHashes)

		printBlock(block, bHash, cHash, c.Id, c.IndexedFields)
	}

	return nil
}

func printBlock(block, bHash, cHash, id string, fields fixity.Fields) {
	fmt.Printf(
		" %s  %s  %s  %s",
		color.GreenString(block),
		color.GreenString(bHash),
		color.YellowString(cHash),
		color.YellowString(id),
	)

	fmt.Print("\n")
}

func sumHash(h string, doNothing bool) string {
	if doNothing {
		return h
	}
	return h[len(h)-8:]
}
