package indexreader

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/storereader"
)

type readerData struct {
	AnchorRand int      `json:"anchorRand"`
	MultiPart  string   `json:"multiPart"`
	Parts      []string `json:"parts"`
	Part       []byte   `json:"part"`
}

type Config struct {
	HashOrAnchor string
	Store        store.Store
	Query        index.Queryer
}

type Reader struct {
	hashOrAnchor string
	store        store.Store
	storeReader  *storereader.Reader
	query        index.Queryer
}

func New(c Config) (*Reader, error) {
	if c.Store == nil {
		return nil, errors.New("missing required config field: Store")
	}
	if c.Query == nil {
		return nil, errors.New("missing required config field: Query")
	}

	return &Reader{
		hashOrAnchor: c.HashOrAnchor,
		store:        c.Store,
		query:        c.Query,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	// if we haven't made a store reader, we need to check if the given hash is
	// an anchor and resolve it to an actual address.
	if r.storeReader == nil {
		h, err := index.ResolveHashOrAnchor(r.store, r.query, r.hashOrAnchor)
		if err != nil {
			return 0, errors.Stack(err)
		}

		sr, err := storereader.New(storereader.Config{
			Hash:  h,
			Store: r.store,
		})
		if err != nil {
			return 0, errors.Stack(err)
		}
		r.storeReader = sr
	}

	return r.storeReader.Read(p)
}
