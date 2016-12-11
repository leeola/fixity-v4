package main

import (
	"io"
	"os"

	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
	"github.com/urfave/cli"
)

func uploadCommand(c *cli.Context) error {
	cType := c.String("content-type")
	if cType == "" {
		return errors.New("missing content type value")
	}

	// Used to automatically attach a filename
	var originalFilename string

	var r io.Reader
	if c.Bool("stdin") {
		r = os.Stdin
		// Continue here, add upload --stdin junk so it uploads files by default
	} else if p := c.Args().First(); p != "" {
		originalFilename = p
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	if r == nil {
		cli.ShowSubcommandHelp(c)
		return errors.New("error: missing path")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	metaChanges := store.MetaChanges{}
	if anchor := c.String("anchor"); anchor != "" {
		metaChanges["anchor"] = anchor
	}
	if newAnchor := c.Bool("new-anchor"); newAnchor != false {
		metaChanges["newAnchor"] = "true"
	}

	if filename := c.String("filename"); filename != "" {
		metaChanges["filename"] = filename
	} else if originalFilename != "" {
		metaChanges["filename"] = originalFilename
	}

	hashes, err := client.Upload(cType, r, metaChanges)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}
