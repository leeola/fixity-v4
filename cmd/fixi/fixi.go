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
			Aliases:   []string{"b"},
			Usage:     "inspect a raw blob from HASH",
			Action:    BlobCmd,
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
			ArgsUsage: "CLIJSON",
			Usage:     "write a commit to fixity",
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
				cli.StringSliceFlag{
					Name:  "index",
					Usage: "a field or field=value to index",
				},
				cli.StringSliceFlag{
					Name:  "fts",
					Usage: "a field or field=value to index with full text search",
				},
				cli.StringFlag{
					Name:  "json-key",
					Usage: "the multijson key to associate the given json with",
					Value: "fixi-cli",
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
