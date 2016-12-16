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

	var r io.Reader = rc
	if !c.Bool("allow-content") {
		teeR, cType, err := getContentType(r)
		if err != nil {
			return err
		}
		r = teeR

		if cType == "Content" {
			Printlnf(`To prevent large content bytes from spamming your console,
viewing blobs of type "Content" is disabled by default.

If you really want to display the full content of this blob,
use the --allow-content flag.
`)
			return errors.New("cannot display content")
		}
	}

	io.Copy(os.Stdout, r)

	return nil
}

func getContentType(r io.Reader) (io.Reader, string, error) {
	var copyR bytes.Buffer

	// Tee the reader so we can check it's type.
	teeR := io.TeeReader(r, &copyR)

	b, err := ioutil.ReadAll(teeR)
	if err != nil {
		return nil, "", err
	}

	var c struct {
		AnchorRand int
		Parts      []string
		Content    []byte

		Anchor string
		Multi  string
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, "", err
	}

	var cType string
	switch {
	case len(c.Content) != 0:
		cType = "Content"
	case len(c.Parts) != 0:
		cType = "MultiPart"
	case c.AnchorRand != 0:
		cType = "Perma"
	case c.Anchor != "" || c.Multi != "":
		cType = "Meta"
	default:
		cType = "Unknown"
	}

	return &copyR, cType, nil
}
