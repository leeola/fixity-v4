package bleve

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/kala"
	kq "github.com/leeola/kala/q"
)

type Config struct {
	// Path is the *directory* containing the bleve index.
	Path string `toml:"path"`
	Log  log15.Logger
}

type Bleve struct {
	config Config
	log    log15.Logger
	bleve  bleve.Index
}

func New(c Config) (*Bleve, error) {
	if c.Path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if err := os.MkdirAll(c.Path, 0755); err != nil {
		return nil, err
	}

	b, err := bleve.Open(c.Path)
	if err == bleve.ErrorIndexMetaMissing {
		mapping := bleve.NewIndexMapping()
		b, err = bleve.New(c.Path, mapping)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to open or create bleve normal index")
	}

	return &Bleve{
		config: c,
		log:    c.Log,
		bleve:  b,
	}, nil
}

func (b *Bleve) Index(h, id string, fields []kala.Field) error {
	b.log.Warn("id indexing not yet implemented", "id", id)

	row := map[string]interface{}{}
	for _, f := range fields {
		// TODO(leeola): implement options
		row[f.Field] = f.Value
	}

	return b.bleve.Index(h, &row)
}

func (b *Bleve) Search(kq *kq.Query) ([]string, error) {
	search, err := ConvertQuery(kq)
	if err != nil {
		return nil, errors.Stack(err)
	}

	searchResults, err := b.bleve.Search(search)
	if err != nil {
		return nil, errors.Stack(err)
	}

	hashes := make([]string, len(searchResults.Hits))
	for i, documentMatch := range searchResults.Hits {
		hashes[i] = documentMatch.ID
	}

	return hashes, nil
}
