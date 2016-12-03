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
			Name:  "config",
			Usage: "use specified kala config",
			Value: "~/.kala.toml",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "id",
			Usage:  "display this nodes id",
			Action: idCommand,
		},
		{
			Name:   "query",
			Usage:  "query the node",
			Action: queryCommand,
		},
	}

	app.Run(os.Args)
}

func Printlnf(f string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(f, v...))
}
