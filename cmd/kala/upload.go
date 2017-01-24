package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/contenttype"
	"github.com/mitchellh/ioprogress"
	"github.com/urfave/cli"
)

func uploadCommand(c *cli.Context) error {
	metaChanges := argsToChanges(c.Args())

	p := c.String("path")
	if !c.Bool("stdin") && p == "" {
		cli.ShowSubcommandHelp(c)
		return errors.New("no upload source given")
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	if c.Bool("stdin") {
		return uploadReader(client, ioutil.NopCloser(os.Stdin), 0, metaChanges)
	}

	fi, err := os.Stat(p)
	if err != nil {
		return err
	}

	// If it's a file, open it and return.
	if !fi.IsDir() {
		return uploadFile(client, p, metaChanges)
	}

	root := p
	return filepath.Walk(p, func(p string, info os.FileInfo, err error) error {
		// If recursive was not enabled and we're beyond the root directory,
		// skip this file.
		if !c.Bool("recursive") && filepath.Clean(root) != filepath.Dir(p) {
			return nil
		}

		// skip hidden unless explicitly included.
		if !c.Bool("hidden") && strings.HasPrefix(filepath.Base(p), ".") {
			return nil
		}

		fi, err := os.Stat(p)
		if err != nil {
			return err
		}
		// Skip all directories outright. Walk takes care of the recursion.
		if fi.IsDir() {
			return nil
		}

		Printlnf("uploading: %s", p)
		return uploadFile(client, p, cloneChanges(metaChanges))
	})
}

func contentTypeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".ogg", ".mov", ".mp4":
		return "video"
	case ".jpg", ".jpeg":
		return "image"
	default:
		return "file"
	}
}

func uploadReader(c *client.Client, rc io.ReadCloser, rSize int64, ch contenttype.Changes) error {
	defer rc.Close()

	// Create the progress reader
	progressR := &ioprogress.Reader{
		Reader:   rc,
		Size:     rSize,
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBytes),
	}

	hashes, err := c.Upload(progressR, ch)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		Printlnf(h)
	}

	return nil
}

func uploadFile(c *client.Client, p string, ch contenttype.Changes) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	// If the user didn't specify metadata to add themselves,
	// automatically set the filename and constenttype for an easy ux.
	if _, ok := ch["filename"]; !ok {
		ch.Set("filename", filepath.Base(p))
	}
	if _, ok := ch.GetContentType(); !ok {
		ch.SetContentType(contentTypeFromExt(filepath.Ext(p)))
	}

	fi, _ := f.Stat()
	return uploadReader(c, f, fi.Size(), ch)
}

func cloneChanges(s contenttype.Changes) contenttype.Changes {
	d := contenttype.Changes{}
	for k, v := range s {
		d[k] = v
	}
	return d
}
