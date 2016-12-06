package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
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
			Name:   "blob",
			Usage:  "get raw blobs",
			Action: blobCommand,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "upload",
				},
			},
		},
		{
			Name:   "id",
			Usage:  "display this nodes id",
			Action: idCommand,
		},
		{
			Name:   "query",
			Action: queryCommand,
		},
		{
			Name:   "upload",
			Action: uploadCommand,
		},
	}

	app.Run(os.Args)
}

func Printlnf(f string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(f, v...))
}
