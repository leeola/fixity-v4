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
			Name:  "config",
			Usage: "use specified kala config",
			Value: "~/.kala.toml",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "query",
			Usage:  "query the index for matching hashes",
			Action: queryCommand,
		},
		{
			Name:      "upload",
			Usage:     "upload content with metadata",
			ArgsUsage: "<FILE>",
			Action:    uploadCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "c, content-type",
					Usage: "the content type of the uploaded data",
					Value: "file",
				},
				cli.StringFlag{
					Name:  "f, filename",
					Usage: "the filename to use in place of the filename on disk",
				},
				cli.StringFlag{
					Name:  "a, anchor",
					Usage: "the anchor to use, if any",
				},
				cli.BoolFlag{
					Name:  "n, new-anchor",
					Usage: "create a new anchor for this upload",
				},
				cli.BoolFlag{
					Name:  "i, stdin",
					Usage: "read from stdin instead of a file",
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
			Name:   "edit",
			Usage:  "download, edit and upload the given hash contents",
			Action: editCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "e, editor",
					Value: "vim",
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
