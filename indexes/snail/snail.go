package snail

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
	"github.com/leeola/fixity/util"
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
	config   Config
	log      log15.Logger
	db       *bolt.DB
	bleveId  bleve.Index
	bleveVer bleve.Index
}

func New(c Config) (*Snail, error) {
	if c.Path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	bleveIdDir := filepath.Join(c.Path, "bleve", "id")
	bleveVerDir := filepath.Join(c.Path, "bleve", "ver")
	// Making the bleve path, because that's one dir deeper than the default,
	// so this'll make everything we need.
	if err := os.MkdirAll(bleveIdDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(bleveVerDir, 0755); err != nil {
		return nil, err
	}

	boltPath := filepath.Join(c.Path, "snail.db")
	db, err := bolt.Open(boltPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	bleveId, err := bleve.Open(bleveIdDir)
	if err == bleve.ErrorIndexMetaMissing {
		mapping := bleve.NewIndexMapping()
		bleveId, err = bleve.New(bleveIdDir, mapping)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to open or create bleve id index")
	}
	bleveVer, err := bleve.Open(bleveVerDir)
	if err == bleve.ErrorIndexMetaMissing {
		mapping := bleve.NewIndexMapping()
		bleveVer, err = bleve.New(bleveVerDir, mapping)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to open or create bleve id index")
	}

	return &Snail{
		config:   c,
		db:       db,
		bleveId:  bleveId,
		bleveVer: bleveVer,
		log:      c.Log,
	}, nil
}

// IndexFts goes through the fields and indexes any with the FTS option.
//
// It ignores all non-fts fields.
func (s *Snail) IndexFts(h, id string, fields []fixity.Field) error {
	docFields := map[string]interface{}{}
	for _, f := range fields {
		opts := f.Options
		if opts == nil {
			continue
		}
		v, ok := opts[fixity.FOKeyFullTextSearch]
		if !ok {
			continue
		}

		useFts, ok := v.(bool)
		if !ok {
			return errors.Errorf("incorrectly option value type: %s", fixity.FOKeyFullTextSearch)
		}

		if useFts {
			docFields[f.Field] = f.Value
		}
	}

	if err := s.bleveVer.Index(h, &docFields); err != nil {
		return err
	}

	if id != "" {
		if err := s.bleveId.Index(h, &docFields); err != nil {
			return err
		}
	}

	return nil
}

// Index the given key by the given document fields.
func (s *Snail) Index(h, id string, fields []fixity.Field) error {
	if err := s.IndexFts(h, id, fields); err != nil {
		return err
	}

	docFields := map[string]interface{}{}
	for _, f := range fields {
		// this is where we'll implement/fork Bleve FTS support.
		// options not supported yet
		if f.Options != nil {
			for optKey, optValue := range f.Options {
				if optKey != fixity.FOKeyFullTextSearch {
					s.log.Warn("snail index option not supported",
						"option", optKey, "value", optValue)
				}
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

			docKey := string(k)
			// nil is okay
			docVal, _ := docFields[q.Constraint.Field]

			// if the doc matches, add it to our doc list to be skipped, limited and
			// sorted.
			if matcher.Match(docKey, q.Constraint.Value, docVal) {
				total += 1

				doc := Doc{
					Key: docKey,
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
	return util.MultiError(
		s.db.Close(),
		s.bleveId.Close(),
		s.bleveVer.Close(),
	)
}
