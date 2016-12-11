package blev

import (
	"github.com/blevesearch/bleve"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
)

type Config struct {
	BleveDir string

	// Optional.
	Log log15.Logger `json:"-"`
}

type Bleve struct {
	index        bleve.Index
	indexVersion string
	log          log15.Logger
}

func New(c Config) (*Bleve, error) {
	if c.BleveDir == "" {
		return nil, errors.New("missing required Config field: BleveDir")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	index, err := bleve.Open(c.BleveDir)

	// if bleve does not exist in the given dir, create it.
	if err == bleve.ErrorIndexMetaMissing {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(c.BleveDir, mapping)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to open or create bleve")
	}

	return &Bleve{
		index: index,
		log:   c.Log,
	}, nil
}
func (b *Bleve) keyExists(k string) (bool, error) {
	doc, err := b.index.Document(k)
	return doc != nil, err
}
