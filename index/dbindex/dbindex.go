package dbindex

import (
	"errors"

	"github.com/leeola/kala/database"
	"github.com/leeola/kala/store"
)

type Config struct {
	Database database.Database
	Store    store.Store
}

type Dbindex struct {
	db    database.Database
	store store.Store
}

func New(c Config) (*Dbindex, error) {
	if c.Database == nil {
		return nil, errors.New("missing required Config field: Database")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}

	return &Dbindex{
		db:    c.Database,
		store: c.Store,
	}, nil
}

func (dbi *Dbindex) AddEntry(h string) error {
	return dbi.db.InsertIndexEntry(h)
}
