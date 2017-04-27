package storm

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/leeola/kala"
	"github.com/leeola/kala/q"
)

// IndexFilename is the filename used as the boltdb index.
const IndexFilename = "index.db"

type Config struct {
	// Path is the *directory* containing the IndexFilename.
	//
	// This will be created if it does not exist.
	Path string `toml:"path"`
}

type Storm struct {
	config Config
	db     *storm.DB
}

func New(c Config) (*Storm, error) {
	if c.Path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if err := os.MkdirAll(c.Path, 0755); err != nil {
		return nil, err
	}

	db, err := storm.Open(filepath.Join(c.Path, IndexFilename))
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
