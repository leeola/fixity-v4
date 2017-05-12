package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/leeola/fixity/autoload"
	"github.com/nwidger/jsoncolor"
	"github.com/urfave/cli"
)

func BlobCmd(ctx *cli.Context) error {
	h := ctx.Args().Get(0)
	if h == "" {
		return cli.ShowCommandHelp(ctx, "blob")
	}

	fixity, err := autoload.LoadFixity(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	b, err := fixity.Blob(h)
	if err != nil {
		return err
	}

	f := jsoncolor.NewFormatter()

	f.SpaceColor = color.New(color.FgRed, color.Bold)
	f.CommaColor = color.New(color.FgWhite, color.Bold)
	f.ColonColor = color.New(color.FgBlue)
	f.ObjectColor = color.New(color.FgBlue, color.Bold)
	f.ArrayColor = color.New(color.FgWhite)
	f.FieldColor = color.New(color.FgGreen)
	f.StringColor = color.New(color.FgBlack, color.Bold)
	f.TrueColor = color.New(color.FgWhite, color.Bold)
	f.FalseColor = color.New(color.FgRed)
	f.NumberColor = color.New(color.FgWhite)
	f.NullColor = color.New(color.FgWhite, color.Bold)

	prettyJson := bytes.Buffer{}
	if err := f.Format(&prettyJson, b); err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, &prettyJson)
	if err != nil {
		return err
	}

	fmt.Print("\n")
	return nil
}
