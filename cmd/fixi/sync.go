package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/leeola/fixity/sync"
	"github.com/urfave/cli"
)

func SyncCmd(ctx *cli.Context) error {
	useCli := ctx.Bool("cli")
	useStdin := ctx.Bool("stdin")
	writingIo := useCli || useStdin

	var path string
	if !writingIo {
		path = ctx.Args().First()
	}

	var cliData string
	if useCli {
		cliData = strings.Join(ctx.Args(), " ")
	}

	var id string
	if writingIo {
		id = ctx.String("id")
	}

	if useCli && useStdin {
		cli.ShowCommandHelp(ctx, "sync")
		return errors.New("error: cannot use --cli and --stdin together")
	}

	if len(ctx.Args()) != 0 && useStdin {
		cli.ShowCommandHelp(ctx, "sync")
		return errors.New("error: cannot use sync path and --stdin together")
	}

	if !writingIo && path == "" {
		cli.ShowCommandHelp(ctx, "sync")
		return errors.New("error: sync path cannot be empty when syncing files")
	}

	if cliData == "" && useCli {
		cli.ShowCommandHelp(ctx, "sync")
		return errors.New("error: must provide cli data with --cli")
	}

	if id == "" && writingIo {
		cli.ShowCommandHelp(ctx, "sync")
		return errors.New("error: --id must be provided with --cli or --stdin")
	}

	fixi, err := loadFixity(ctx)
	if err != nil {
		return err
	}

	fields, err := fieldsFromCtx(ctx)
	if err != nil {
		return err
	}

	if writingIo {
		var r io.Reader
		if useCli {
			r = strings.NewReader(cliData)
		} else if useStdin {
			r = os.Stdin
		}
		s := sync.Io(fixi, id, r, os.Stdout, fields...)

		for more := s.Next(); more; more = s.Next() {
			// no printing needed, io writes to the given stdout.
			if _, err := s.Value(); err != nil {
				return err
			}
		}

		return nil
	}

	// syncConf := sync.Config{
	// 	Path:      path,
	// 	Folder:    ctx.String("folder"),
	// 	Recursive: ctx.Bool("recursive"),
	// 	Fixity:    fixi,
	// }
	// sync, err := sync.New(syncConf)
	// if err != nil {
	// 	return err
	// }

	// go func() {
	// 	for msg := range sync.Updates() {
	// 		fmt.Println(msg)
	// 	}
	// }()

	// return sync.Sync()

	return nil
}
