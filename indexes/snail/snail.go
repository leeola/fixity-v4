package snail

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
)

type DocFields map[string]interface{}

type Doc struct {
	Key    string
	Fields DocFields
}

type Config struct {
	// Path is the *directory* containing the snail index.
	Path string `toml:"path"`

	// Bucket name for the root bucket
	Bucket string `toml:"bucket"`

	Log log15.Logger
}

// Snail is a slow implementation of indexing fields to match given keys.
//
// Snail is works by indexing the values of all fields in a sorted manner,
// and then upon querying, it filters the resulting keys with each constraint.
// Duplicate keys will be replaced when present.
type Snail struct {
	config Config
	log    log15.Logger
	db     *bolt.DB
}

func New(c Config) (*Snail, error) {
	if c.Path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if err := os.MkdirAll(c.Path, 0755); err != nil {
		return nil, err
	}

	boltPath := filepath.Join(c.Path, "snail.db")
	db, err := bolt.Open(boltPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Snail{
		config: c,
		db:     db,
		log:    c.Log,
	}, nil
}

// Index the given key by the given document fields.
func (s *Snail) Index(h, id string, fields []fixity.Field) error {
	docFields := map[string]interface{}{}
	for _, f := range fields {
		// this is where we'll implement/fork Bleve FTS support.
		// options not supported yet
		if f.Options != nil {
			for optKey, optValue := range f.Options {
				s.log.Warn("snail index option not implemented",
					"option", optKey, "value", optValue)
			}
		}
		docFields[f.Field] = f.Value
	}

	b, err := json.Marshal(docFields)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		verBkt, err := tx.CreateBucketIfNotExists([]byte("version"))
		if err != nil {
			return errors.Errorf("create bucket: %s", err)
		}

		var idBkt *bolt.Bucket
		if id != "" {
			bkt, err := tx.CreateBucketIfNotExists([]byte("id"))
			if err != nil {
				return errors.Errorf("create bucket: %s", err)
			}
			idBkt = bkt
		}

		// TODO(leeola): compare speed of keying documents by hash vs by
		// integers. This isn't a big concern right now, but documents are being
		// keyed by hash, so that repeated indexes of the same content will replace
		// old indexes and not duplicate documents.
		//
		// Note that this is optional for hashes, but required for id indexing.
		if err := verBkt.Put([]byte(h), b); err != nil {
			return err
		}

		if id != "" {
			if err := idBkt.Put([]byte(h), b); err != nil {
				return err
			}
		}

		return nil
	})
}

// TODO(leeola): I'll likely break Snail off into it's own package, which
// means packaging it's own set of query constraints/ops/etc. For now though,
// i'm piggybacking Fixity's set of constraints.
func (s *Snail) Search(q *q.Query) ([]string, error) {
	var matchedDocs []Doc

	// TODO(leeola): Make this configurable, once the Version vs Id searches
	// are implemented into fixity.
	bktName := []byte("version")

	matcher, err := s.convertConstraint(q.Constraint)
	if err != nil {
		return nil, err
	}

	hasSorts := len(q.SortBy) > 0
	total := 0

	err = s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktName)
		// if the bucket does not exist, no data has been indexed.
		// This is not an error, just that no results will be found.
		if bkt == nil {
			return nil
		}

		c := bkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var docFields DocFields

			if err := json.Unmarshal(v, &docFields); err != nil {
				return err
			}

			// if the doc matches, add it to our doc list to be skipped, limited and
			// sorted.
			if matcher(q.Constraint, docFields) {
				total += 1

				doc := Doc{
					Key: string(k),
				}

				// only store the if we need to sort. This helps reduce memory footprint if
				// we don't actually need to store all the documents.
				//
				// TODO(leeola): only store the fields we're sorting by, further reducing
				// the memory footprint.
				if hasSorts {
					doc.Fields = docFields
				}

				matchedDocs = append(matchedDocs, doc)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// TODO(leeola): Sort the dataset based on each doc field.

	// Now that the docs are sorted, skip/paginate as needed.

	if q.SkipBy >= total {
		return nil, nil
	}

	endIndex := q.SkipBy + q.LimitBy
	if endIndex > total {
		endIndex = total
	}

	// apply the skip and limit
	matchedDocs = matchedDocs[q.SkipBy:endIndex]

	// The actual keys for the matchedDoc
	matches := make([]string, endIndex-q.SkipBy)
	for i, d := range matchedDocs {
		matches[i] = d.Key
	}

	return matches, nil
}

func (s *Snail) Close() error {
	return s.db.Close()
}
