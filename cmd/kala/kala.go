package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "a kala store cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "kala host",
			Value: "",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "use specified kala config",
			Value:  "~/.kala.toml",
			EnvVar: "KALA_CONFIG",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "query",
			Usage:  "query the index for matching hashes",
			Action: queryCommand,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "a, ascending",
					Usage: "a repeatable ascending meta `FIELD` to sort by (default: [indexEntry])",
				},
				cli.StringSliceFlag{
					Name:  "d, descending",
					Usage: "a repeatable descending meta `FIELD` to sort by",
				},
				cli.StringSliceFlag{
					Name:  "p, preview",
					Usage: "include meta `FIELD` in the query results",
				},
			},
		},
		{
			Name:      "upload",
			Usage:     "upload content with metadata",
			ArgsUsage: "<FILE>",
			Action:    uploadCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "p, path",
					Usage: "a file or directory to upload from disk",
				},
				cli.BoolFlag{
					Name:  "i, stdin",
					Usage: "read from stdin instead of a file",
				},
				cli.BoolFlag{
					Name:  "r, recursive",
					Usage: "upload files from dirs recursively",
				},
				cli.BoolFlag{
					Name:  "hidden",
					Usage: "include hidden files",
				},
			},
		},
		{
			Name:      "download",
			Action:    downloadCommand,
			Usage:     "download content with metadata",
			ArgsUsage: "<FILE>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "stdout",
				},
			},
		},
		{
			Name:   "meta",
			Usage:  "change metadata for the given anchor",
			Action: metaCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, anchor",
					Usage: "specify an anchor to base the metadata change off of",
				},
				cli.StringFlag{
					Name:  "m, meta",
					Usage: "specify a meta to base the metadata change off of",
				},
			},
		},
		{
			Name:      "blob",
			Usage:     "print a blob hash",
			ArgsUsage: "<hash>",
			Action:    blobCommand,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "a, allow-content",
				},
			},
		},
		{
			Name:   "id",
			Usage:  "print the id of the connected node",
			Action: idCommand,
		},
	}

	app.Run(os.Args)
}

func Printlnf(f string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(f, v...))
}
