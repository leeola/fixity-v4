package main

import (
	"fmt"
	"os"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/autoload"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fixi"
	app.HelpName = "fixi" // this was being set to "blob", how?
	app.Usage = "a low level cli to interact with a fixity store"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "~/.config/fixity/config.toml",
			Usage:  "load config from `PATH`",
			EnvVar: "FIXI_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "blob",
			ArgsUsage: "HASH",
			Usage:     "inspect a raw blob from HASH",
			Action:    BlobCmd,
		},
		{
			Name:  "blocks",
			Usage: "inspect the fixity blockchain",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "block-hashes",
					Usage: "display full block hashes",
				},
				cli.BoolFlag{
					Name:  "content-hashes",
					Usage: "display full content hashes",
				},
				cli.IntFlag{
					Name:  "limit",
					Value: 25,
					Usage: "limit the total blocks displayed by `LIMIT`",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "only display blocks of `TYPE`",
				},
			},
			Action: BlocksCmd,
		},
		{
			Name:      "read",
			ArgsUsage: "ID PATH",
			Aliases:   []string{"r"},
			Usage:     "read the given ID or hash content to the given PATH",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "hash",
					Usage: "use a hash as the first arg instead of hash",
				},
				cli.BoolFlag{
					Name:  "stdout",
					Usage: "print the bytes to stdout; this can be verbose!",
				},
			},
			Action: ReadCmd,
		},
		{
			Name:      "search",
			ArgsUsage: "QUERY",
			Aliases:   []string{"s"},
			Usage:     "search for hashes matching the query",
			Action:    SearchCmd,
		},
		{
			Name:      "write",
			Aliases:   []string{"w"},
			ArgsUsage: "FILE",
			Usage:     "write a content to fixity",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "cli",
					Usage: "upload from cli args",
				},
				cli.BoolFlag{
					Name:  "stdin",
					Usage: "upload from stdin",
				},
				cli.BoolFlag{
					Name:  "spam-bytes",
					Usage: "do not hide large chunks byte contents",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "the id of the content",
				},
				cli.IntFlag{
					Name:  "manual-rollsize",
					Usage: "the rollsize in B",
				},
				cli.StringSliceFlag{
					Name:  "index",
					Usage: "a field or field=value to index",
				},
				cli.StringSliceFlag{
					Name:  "fts",
					Usage: "a field or field=value to index with full text search",
				},
				cli.BoolFlag{
					Name:  "inspect",
					Usage: "inspect the written data structure",
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

// loadFixity expands the configPath and loads fixity.
func loadFixity(ctx *cli.Context) (fixity.Fixity, error) {
	configPath := ctx.GlobalString("config")

	configPath, err := homedir.Expand(configPath)
	if err != nil {
		return nil, err
	}

	return autoload.LoadFixity(configPath)
}
