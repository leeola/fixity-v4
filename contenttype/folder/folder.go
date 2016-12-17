package folder

import (
	"errors"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type FolderMeta struct {
	store.Meta

	Foldername string `json:"foldername"`
}

type Config struct {
	Store store.Store
	Index index.Indexer
}

type Folder struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*Folder, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &Folder{
		store: c.Store,
		index: c.Index,
	}, nil
}
