package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/leeola/kala/client"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/urfave/cli"
)

// ImportExport stores a generic piece of data with metadata in the form of changes.
//
// This exists to allow for exporting and importing generic Kala data without
// having to know the contenttype schema/format. ContentPath is just the raw data,
// eg an image, and changes is a map of changes as would be passed in during
// the initial upload.
//
// While this may seem like an odd way to "export" data that has a sane disk
// representation, like images, this also is intended to handle data that does
// not have a disk representation like tweets, facebook posts, etc. That type
// of data would require special representations to sit on disk. Afterwards
// importing that data back into a Kala store would require different handling
// based on how natural it's disk representation is.
//
// To avoid this, everything is treated the same way. Data, and KeyValue Changes
// to be applied during upload just as if you were uploading the data and setting
// `title=foo tags=bar` by hand.
type ImportExport struct {
	// Changes is a map of contenttype key=value changes, such as title=foo.
	//
	// Note that how changes is handled is up to the ContentType implementor.
	Meta contenttype.Changes `json:"meta"`

	// ContentPath is the relative path to the content for the given changes.
	//
	// Note that this is optional, a piece of metadata may not have any content.
	ContentPath string `json:"contentPath,omitempty"`
}

func importCommand(ctx *cli.Context) error {
	srcDir := ctx.Args().Get(0)
	if srcDir == "" {
		cli.ShowSubcommandHelp(ctx)
		return errors.New("missing import path")
	}

	fi, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		cli.ShowSubcommandHelp(ctx)
		return errors.New("import path must be a directory")
	}

	c, err := ClientFromContext(ctx)
	if err != nil {
		return err
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".meta" {
			return nil
		}

		return importHash(c, path)
	})
	if err != nil {
		return err
	}

	return nil
}

func importHash(c *client.Client, p string) error {
	metaF, err := os.OpenFile(p, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer metaF.Close()

	var ie ImportExport
	if err := jsonutil.UnmarshalReader(metaF, &ie); err != nil {
		return err
	}

	// always use a new anchor when importing.
	ie.Meta.Set("newAnchor", "true")

	// If the import has a content path, upload the full content. Otherwise,
	// only upload meta.
	var h string
	if ie.ContentPath != "" {
		contentF, err := os.OpenFile(p, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer contentF.Close()

		hs, err := c.Upload(contentF, ie.Meta)
		if err != nil {
			return err
		}
		h = hs[len(hs)-1]
	} else {
		hs, err := c.UploadMeta(ie.Meta)
		if err != nil {
			return err
		}
		h = hs[len(hs)-1]
	}

	cType, _ := ie.Meta.GetContentType()
	Printlnf("imported %s: %s", cType, h)

	return nil
}
