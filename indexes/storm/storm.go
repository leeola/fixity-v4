package storm

import (
	"errors"

	"github.com/asdine/storm"
	"github.com/leeola/kala"
	"github.com/leeola/kala/q"
)

type Config struct {
	DbPath string `toml:"dbPath"`
}

type Storm struct {
	config Config
	db     *storm.DB
}

func New(c Config) (*Storm, error) {
	db, err := storm.Open(c.DbPath)
	if err != nil {
		return nil, err
	}

	return &Storm{
		config: c,
		db:     db,
	}, nil
}

func (s *Storm) Index(fields []kala.Field) error {
	row := map[string]interface{}{}
	for _, f := range fields {
		// TODO(leeola): implement options
		row[f.Field] = f.Value
	}

	return s.db.Save(row)
}

func (s *Storm) Search(q.Query) ([]string, error) {
	return nil, errors.New("not implemented")
}
