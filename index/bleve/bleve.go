package bleve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/leeola/fixity/config"
)

const (
	subDir      = "index_bleve"
	idIndexDir  = "id_index"
	refIndexDir = "ref_index"
)

type Config struct {
	Path string `json:"path"`
}

type Index struct {
	idIndex  bleve.Index
	refIndex bleve.Index
}

func New(name string, cfg config.Config) (*Index, error) {
	var c Config
	if err := cfg.IndexConfig(name, &c); err != nil {
		return nil, fmt.Errorf("indexconfig: %v", err)
	}

	if c.Path == "" {
		return nil, fmt.Errorf("missing required config field: path")
	}

	idPath := filepath.Join(c.Path, subDir, idIndexDir)
	refPath := filepath.Join(c.Path, subDir, refIndexDir)

	idIndex, err := newBleve(idPath)
	if err != nil {
		return nil, fmt.Errorf("newBleve: %v", err)
	}

	refIndex, err := newBleve(refPath)
	if err != nil {
		return nil, fmt.Errorf("newBleve: %v", err)
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
		index, err = bleve.New(path, newMapping())
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

func newMapping() *mapping.IndexMappingImpl {
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	indexMapping := bleve.NewIndexMapping()

	// ids with non-alpha-num values were having trouble matching,
	// such as "foo-bar". After searching, it appears a keyword
	// analyzer is needed to allow the field to not be chopped up.
	//
	// ref: https://github.com/blevesearch/bleve/issues/844
	indexMapping.DefaultMapping.AddFieldMappingsAt(fieldNameID, keywordFieldMapping)
	indexMapping.DefaultMapping.AddFieldMappingsAt(fieldNameRef, keywordFieldMapping)

	return indexMapping
}
