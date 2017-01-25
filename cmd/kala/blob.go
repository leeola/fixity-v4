package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"github.com/urfave/cli"
)

func blobCommand(c *cli.Context) error {
	h := c.Args().Get(0)

	if h == "" {
		return cli.ShowCommandHelp(c, "blob")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	rc, err := client.GetBlob(h)
	if err != nil {
		return err
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	if !c.Bool("allow-content") {
		cType, err := getContentType(b)
		if err != nil {
			return err
		}

		if cType == "Part" {
			Printlnf(`To prevent large content bytes from spamming your console,
viewing blobs of type "Content" is disabled by default.

If you really want to display the full content of this blob,
use the --allow-content flag.
`)
			return errors.New("cannot display content")
		}
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
	return err
}

func getContentType(b []byte) (string, error) {
	var c struct {
		Meta   string
		Parts  []string
		Part   []byte
		Hashes []string
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return "", err
	}

	var cType string
	switch {
	case c.Meta != "":
		cType = "Version"
	case len(c.Part) != 0:
		cType = "Part"
	case len(c.Parts) != 0:
		cType = "MultiPart"
	case len(c.Hashes) != 0:
		cType = "MultiHash"
	default:
		cType = "Unknown"
	}

	return cType, nil
}
