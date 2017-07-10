package main

import (
	"fmt"

	"github.com/leeola/fixity/sync"
	"github.com/urfave/cli"
)

func SyncCmd(ctx *cli.Context) error {
	path := ctx.Args().First()
	if path == "" {
		return cli.ShowCommandHelp(ctx, "sync")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	syncConf := sync.Config{
		Path:      path,
		Folder:    ctx.String("folder"),
		Recursive: ctx.Bool("recursive"),
		Fixity:    fixi,
	}
	sync, err := sync.New(syncConf)
	if err != nil {
		return err
	}

	go func() {
		for msg := range sync.Updates() {
			fmt.Println(msg)
		}
	}()

	return sync.Sync()
}
