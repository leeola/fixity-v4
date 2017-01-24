package data

import (
	"io"

	"github.com/leeola/errors"
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/ctutil"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
)

type Config struct {
	Store store.Store
	Index index.Indexer
}

type Data struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*Data, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &Data{
		store: c.Store,
		index: c.Index,
	}, nil
}

// TODO(leeola): centralize the common tasks in this method into helpers.
// A lot of this (writing content roller and multipart, etc) is going to be
// duplicated on every ContentType handler.
func (f *Data) StoreContent(rc io.ReadCloser, v store.Version, c ct.Changes) ([]string, error) {
	if rc == nil {
		return nil, errors.New("missing ReadCloser")
	}
	defer rc.Close()

	roller, err := camli.New(rc)
	if err != nil {
		return nil, errors.Stack(err)
	}

	// write the actual content
	hashes, err := ctutil.WritePartRoller(f.store, f.index, roller)
	if err != nil {
		return nil, errors.Stack(err)
	}

	// write the multipart
	h, err := store.WriteMultiPart(f.store, store.MultiPart{
		Parts: hashes,
	})
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)
	c.SetMultiHash(h)

	// Write the entries, not including the final metadata hash
	// The last hash is metadata, and we'll add that manually.
	for _, h := range hashes {
		if err := f.index.Entry(h); err != nil {
			return nil, errors.Stack(err)
		}
	}

	metaHashes, err := f.StoreMeta(v, c)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return append(hashes, metaHashes...), nil
}

func (f *Data) StoreMeta(v store.Version, c ct.Changes) ([]string, error) {
	var meta ct.Meta

	if v.Meta != "" {
		if err := store.ReadAndUnmarshal(f.store, v.Meta, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	mH, vH, err := ct.WriteMetaAndVersion(f.store, f.index, v, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return []string{mH, vH}, nil
}

func (d *Data) MetaToChanges([]byte) (ct.Changes, error) {
	return nil, errors.New("not implemented")
}
