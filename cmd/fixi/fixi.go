package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fixi"
	app.HelpName = "fixi" // this was being set to "blob", how?
	app.Usage = "interact with your fixi datastore"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, f",
			Value:  "~/.config/fixi/fixi.toml",
			Usage:  "load config from `PATH`",
			EnvVar: "FIXI_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "blob",
			ArgsUsage: "HASH",
			Aliases:   []string{"b"},
			Usage:     "inspect a raw blob from HASH",
			Action:    BlobCmd,
		},
		{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "write a commit to fixity",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "upload a blob from `PATH`",
				},
				cli.StringFlag{
					Name:  "blob",
					Usage: "upload a blob from `DATA`",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "the id of the commit",
				},
				cli.StringFlag{
					Name:  "previous",
					Usage: "the previousVersionHash of the commit",
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
