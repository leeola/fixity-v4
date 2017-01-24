package data

import (
	"io"

	"github.com/leeola/errors"
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
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

func (d *Data) StoreContent(rc io.ReadCloser, v ct.Version, c ct.Changes) ([]string, error) {
	h, hashes, err := ct.WriteContent(d.store, d.index, rc)
	if err != nil {
		return nil, errors.Stack(err)
	}
	c.SetMultiHash(h)

	metaHashes, err := d.StoreMeta(v, c)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return append(hashes, metaHashes...), nil
}

func (d *Data) StoreMeta(v ct.Version, c ct.Changes) ([]string, error) {
	var meta ct.Meta

	// Apply any changes to the version, as needed.
	v.FromChanges(c)

	if v.Meta != "" {
		if err := store.ReadAndUnmarshal(d.store, v.Meta, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// update the meta with any changes matching the meta
	meta.FromChanges(c)

	mH, vH, err := ct.WriteMetaAndVersion(d.store, d.index, v, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return []string{mH, vH}, nil
}

func (d *Data) MetaToChanges([]byte) (ct.Changes, error) {
	return nil, errors.New("not implemented")
}
