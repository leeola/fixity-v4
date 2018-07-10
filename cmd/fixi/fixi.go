package main

import (
	"fmt"
	"os"

	"github.com/leeola/fixity/blobstore/disk"
	"github.com/leeola/fixity/index/bleve"
	"github.com/leeola/fixity/store/nosign"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fixi"
	app.Usage = "a low level cli to interact with a fixity store"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "~/.config/fixity/client.toml",
			Usage:  "load config from `PATH`",
			EnvVar: "FIXI_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "blob",
			ArgsUsage: "HASH",
			Usage:     "inspect a blob from HASH",
			Action:    BlobCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "allow-unsafe",
					Usage: "allow printing schemaless bytes",
				},
			},
		},
		{
			Name:      "read",
			ArgsUsage: "HASH",
			Usage:     "read a mutation from HASH",
			Action:    ReadCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dont-print-values",
					Usage: "do not output values to stderr",
				},
				cli.BoolFlag{
					Name:  "no-stderr-color",
					Usage: "do not output color to stderr",
				},
				cli.BoolFlag{
					Name:  "no-mutation",
					Usage: "do not print mutation to stderr",
				},
				cli.BoolFlag{
					Name:  "no-values",
					Usage: "do not print values to stderr",
				},
				cli.StringFlag{
					Name:  "filename",
					Usage: "output data to given filename",
				},
			},
		},
		{
			Name:      "write",
			Aliases:   []string{"w"},
			ArgsUsage: "FILE",
			Usage:     "write a content to fixity",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "id of written data",
				},
				cli.BoolFlag{
					Name:  "stdin",
					Usage: "upload from stdin",
				},
				cli.BoolFlag{
					Name:  "preview",
					Usage: "preview blobs with schemas",
				},
				cli.BoolFlag{
					Name:  "allow-unsafe",
					Usage: "allow previewing schemaless bytes",
				},
			},
			Action: WriteCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func storeFromCli(clictx *cli.Context) (*nosign.Store, error) {
	path := "./_store"

	bs, err := disk.New(path, true)
	if err != nil {
		return nil, fmt.Errorf("blobstore new: %v", err)
	}

	ix, err := bleve.New(path)
	if err != nil {
		return nil, fmt.Errorf("index new: %v", err)
	}

	s, err := nosign.New(bs, ix)
	if err != nil {
		return nil, fmt.Errorf("store new: %v", err)
	}

	return s, nil
}
