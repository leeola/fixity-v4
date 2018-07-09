package bleve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
)

const (
	subDir      = "index_bleve"
	idIndexDir  = "id_index"
	refIndexDir = "ref_index"
)

type Index struct {
	idIndex  bleve.Index
	refIndex bleve.Index
}

func New(path string) (*Index, error) {
	idPath := filepath.Join(path, subDir, idIndexDir)
	refPath := filepath.Join(path, subDir, refIndexDir)

	idIndex, err := newBleve(idPath)
	if err != nil {
		return nil, fmt.Errorf("newBleve:  %v", err)
	}

	refIndex, err := newBleve(refPath)
	if err != nil {
		return nil, fmt.Errorf("newBleve:  %v", err)
	}

	return &Index{
		idIndex:  idIndex,
		refIndex: refIndex,
	}, nil
}

func newBleve(path string) (bleve.Index, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("mkdirall %s: %v", path, err)
	}

	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexMetaMissing {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(path, mapping)
		if err != nil {
			return nil, fmt.Errorf("new ref index: %v", err)
		}
		return index, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open ref index: %v", err)
	}
	return index, nil
}
