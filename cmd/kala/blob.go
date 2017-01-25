package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

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

	var buf bytes.Buffer
	if err := json.Indent(&buf, b, "", "\t"); err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, &buf)
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
