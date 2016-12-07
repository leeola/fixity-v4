package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

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
		var copyR bytes.Buffer
		r = &copyR

		// Tee the reader so we can check it's type.
		teeR := io.TeeReader(rc, &copyR)
		cType, err := getContentType(teeR)
		if err != nil {
			return err
		}

		if cType == "Content" || cType == "Unknown" {
			Printlnf(`To prevent large content bytes from spamming your console,
viewing blobs of type "Content" or "Unknown" is disabled by default.

If you really want to display the full content of this blob,
use the --allow-content flag.
`)
			return errors.New("cannot display content")
		}
	}

	io.Copy(os.Stdout, r)

	return nil
}

func getContentType(r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	var c struct {
		CreatedAt time.Time
		Rand      int
		Parts     []string
		Content   []byte
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return "", err
	}

	switch {
	case len(c.Content) != 0:
		return "Content", nil
	case len(c.Parts) != 0:
		return "MultiPart", nil
	case !c.CreatedAt.IsZero() && c.Rand != 0:
		return "Perma", nil
	default:
		return "Unknown", nil
	}
}
