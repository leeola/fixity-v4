package bleve

import (
	"fmt"
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
	path = filepath.Join(path, subDir)

	// open the unique index. If it does not exist, it will be created under Rebuild
	idIndex, err := bleve.Open(filepath.Join(path, idIndexDir))
	if err != nil && err != bleve.ErrorIndexMetaMissing {
		return nil, fmt.Errorf("open id index: %v", err)
	}

	refIndex, err := bleve.Open(filepath.Join(path, refIndexDir))
	if err != nil && err != bleve.ErrorIndexMetaMissing {
		return nil, fmt.Errorf("open ref index: %v", err)
	}

	return &Index{
		idIndex:  idIndex,
		refIndex: refIndex,
	}, nil
}
