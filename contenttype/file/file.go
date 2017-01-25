package file

import (
	"io"

	"github.com/leeola/errors"
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

// The key used to address this contenttype by in maps.
const TypeKey = "file"

type Config struct {
	Store store.Store
	Index index.Indexer
}

type File struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*File, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &File{
		store: c.Store,
		index: c.Index,
	}, nil
}

func (d *File) StoreContent(rc io.ReadCloser, v ct.Version, c ct.Changes) ([]string, error) {
	h, hashes, err := ct.WriteContent(d.store, d.index, rc)
	if err != nil {
		return nil, errors.Stack(err)
	}
	c.SetMultiPart(h)

	metaHashes, err := d.StoreMeta(v, c)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return append(hashes, metaHashes...), nil
}

func (d *File) StoreMeta(v ct.Version, c ct.Changes) ([]string, error) {
	var meta Meta

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

func (f *File) MetaToChanges(b []byte) (ct.Changes, error) {
	return nil, errors.New("not implemented")
	// m := FileMeta{}
	// if len(b) != 0 {
	// 	if err := json.Unmarshal(b, &m); err != nil {
	// 		return nil, errors.Stack(err)
	// 	}
	// }

	// c := contenttype.Changes{}
	// m.ToChanges(c)
	// return c, nil
}

// func UnmarshalMetadata(b []byte) (index.Indexable, error) {
// 	m := FileMeta{}
//
// 	if len(b) != 0 {
// 		if err := json.Unmarshal(b, &m); err != nil {
// 			return nil, errors.Stack(err)
// 		}
// 	}
//
// 	return m, nil
// }
