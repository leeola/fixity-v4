package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/util"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/urfave/cli"
)

func exportCommand(c *cli.Context) error {
	destDir := c.String("dest")

	if destDir == "" {
		destDir = strings.ToLower(
			"kala-export_" + time.Now().Format("Jan-02-2006_15-04-05"))

		Printlnf("exporting to: %s", destDir)
	}

	// Stat the dir to make sure we can write to it.
	_, err := os.Stat(destDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return errors.New("error: destination must not yet exist")
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	client, err := ClientFromContext(c)
	if err != nil {
		return err
	}

	q := index.Query{
		FromEntry:      1,
		SearchVersions: false,
	}

	for h := range client.QueryChan(nil, q) {
		if h.Error != nil {
			return h.Error
		}

		if err := exportHash(client, destDir, h.Hash.Hash); err != nil {
			return err
		}
	}

	Printlnf("exporting finished successfully")
	Printlnf("exported to: %s", destDir)
	return nil
}

func exportHash(c *client.Client, destDir, h string) error {
	changes, err := c.GetDownloadMetaExport(h)
	if err != nil {
		return err
	}

	cType, _ := c.GetBlobContentType(h)
	if cType == "" {
		cType = "_unknown_contenttype"
	}

	destDir = filepath.Join(destDir, cType)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	Printlnf("exporting %s: %s", cType, h)

	var contentPath string
	// TODO(leeola): unset contentPath if the metahash doesn't have any content.
	rc, err := c.Download(h)
	if err != nil {
		return err
	}
	defer rc.Close()
	// For a humanized ux, include the filename from the content, if it exists.
	if f, ok := changes.GetString("filename"); ok {
		contentPath = filepath.Join(destDir, fmt.Sprintf("%s.%s.content", h, f))
	} else {
		contentPath = filepath.Join(destDir, fmt.Sprintf("%s.content", h))
	}
	if err := util.WriteReader(rc, contentPath, 0644); err != nil {
		return err
	}

	metaPath := filepath.Join(destDir, fmt.Sprintf("%s.meta", h))
	err = jsonutil.MarshalToPath(metaPath, ImportExport{
		Meta:        changes,
		ContentPath: contentPath,
	})
	if err != nil {
		return err
	}

	return nil
}
