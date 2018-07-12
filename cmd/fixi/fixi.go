package main

import (
	"fmt"
	"os"

	// import defaults
	_ "github.com/leeola/fixity/defaultpkg"

	"github.com/leeola/fixity"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fixi"
	app.Usage = "a low level cli to interact with a fixity store"
	app.Flags = []cli.Flag{
	// cli.StringFlag{
	// 	Name:   "config, c",
	// 	Value:  "~/.config/fixity/client.toml",
	// 	Usage:  "load config from `PATH`",
	// 	EnvVar: "FIXI_CONFIG",
	// },
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
			Name:      "cat",
			Aliases:   []string{"r"},
			ArgsUsage: "ID",
			Usage:     "read a mutation from ID",
			Action:    CatCmd,
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
				cli.BoolFlag{
					Name:  "ref",
					Usage: "read from mutation refs, not ids",
				},
			},
		},
		{
			Name:      "query",
			Aliases:   []string{"q"},
			ArgsUsage: "QUERY",
			Usage:     "search the store for QUERY",
			Action:    QueryCmd,
			Flags:     []cli.Flag{},
		},
		{
			Name:      "read",
			Aliases:   []string{"r"},
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
				cli.BoolFlag{
					Name:  "ref",
					Usage: "read from mutation refs, not ids",
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
				cli.StringSliceFlag{
					Name:  "kv",
					Usage: "a key=value pair to index write",
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

func storeFromCli(clictx *cli.Context) (fixity.Store, error) {
	return fixity.New()
}
