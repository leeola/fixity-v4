package indexreader

import (
	"bytes"
	"encoding/json"
	"io"

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
	Hash  string
	Store store.Store
	Query index.Queryer
}

type Reader struct {
	storeReader *storereader.Reader
	query       index.Queryer
}

func New(c Config) (*Reader, error) {
	if c.Query == nil {
		return nil, errors.New("missing required config field: Query")
	}

	r, err := storereader.New(storereader.Config{
		Hash:  c.Hash,
		Store: c.Store,
	})
	if err != nil {
		return nil, errors.Stack(err)
	}

	return &Reader{
		storeReader: r,
		query:       c.Query,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	n, hwb, err := r.storeReader.ReadContentOnly(p)
	if err == io.EOF {
		return 0, io.EOF
	}
	if err != nil {
		return 0, errors.Stack(err)
	}

	// If no bytes are returned, this is a normal Reader cycle with some (or zero) data
	// being read into p. Return n.
	if hwb.Bytes == nil {
		return n, nil
	}

	var d readerData
	if err := json.Unmarshal(hwb.Bytes, &d); err != nil {
		return 0, errors.Stack(err)
	}

	switch {
	case d.AnchorRand != 0:
		q := index.Query{
			Metadata: index.Metadata{
				"anchor": hwb.Hash,
			},
		}
		s := index.SortBy{
			Field:      "uploadedAt",
			Descending: true,
		}

		result, err := r.query.QueryOne(q, s)
		if err != nil {
			return 0, errors.Wrap(err, "failed to query anchor")
		}

		if result.Hash.Hash == "" {
			return 0, errors.New("index Reader: hash not found")
		}

		r.storeReader.AddHashes(result.Hash.Hash)

	case d.MultiPart != "":
		r.storeReader.AddHashes(d.MultiPart)

	case len(d.Parts) > 0:
		r.storeReader.AddHashes(d.Parts...)

	case len(d.Part) > 0:
		r.storeReader.SetCurrentReader(bytes.NewReader(d.Part))

	default:
		return 0, errors.New("Reader: unhandled hash content")
	}

	return 0, nil
}
