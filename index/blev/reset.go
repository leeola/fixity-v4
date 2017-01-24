package blev

import (
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/leeola/errors"
)

const (
	anchorIndexDir = "anchor"
	entryIndexDir  = "entries"
)

func (b *Bleve) Reset() error {
	b.log.Debug("resetting bleve index")

	if err := removeIndexDirs(b.indexDir); err != nil {
		return errors.Stack(err)
	}

	if err := makeIndexDirs(b.indexDir); err != nil {
		return errors.Stack(err)
	}

	iv, err := NewVersion(b.indexDir)
	if err != nil {
		return errors.Stack(err)
	}
	b.indexVersion = iv

	mapping := bleve.NewIndexMapping()
	entryIndex, err := bleve.New(filepath.Join(b.indexDir, entryIndexDir), mapping)
	if err != nil {
		return errors.Wrap(err, "failed to create entry index from mapping")
	}
	b.entryIndex = entryIndex

	mapping = bleve.NewIndexMapping()
	anchorIndex, err := bleve.New(filepath.Join(b.indexDir, anchorIndexDir), mapping)
	if err != nil {
		return errors.Wrap(err, "failed to create anchor index from mapping")
	}
	b.anchorIndex = anchorIndex

	return nil
}
