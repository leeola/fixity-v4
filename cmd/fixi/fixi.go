package main

import (
	"fmt"
	"os"

	"github.com/leeola/fixity/blobstore/disk"
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
					Name:  "not-safe",
					Usage: "do not prevent printing schemaless bytes",
				},
			},
		},
		{
			Name:      "write",
			Aliases:   []string{"w"},
			ArgsUsage: "FILE",
			Usage:     "write a content to fixity",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stdin",
					Usage: "upload from stdin",
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
	bs, err := disk.New("./_store")
	if err != nil {
		return nil, fmt.Errorf("blobstore new: %v", err)
	}

	s, err := nosign.New(bs)
	if err != nil {
		return nil, fmt.Errorf("store new: %v", err)
	}

	return s, nil
}
