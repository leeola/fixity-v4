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
					Name:  "type",
					Usage: "standard unix file metadata",
					Value: "file",
				},
				cli.BoolFlag{
					Name:  "stdin",
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
			Name:      "blob",
			Usage:     "print a blob hash",
			ArgsUsage: "<hash>",
			Action:    blobCommand,
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
