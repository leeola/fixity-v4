package inventory

import (
	"encoding/json"
	"io"

	"github.com/leeola/errors"
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

// The key used to address this contenttype by in maps.
const TypeKey = "inventory"

type Config struct {
	Store store.Store
	Index index.Indexer
}

type Inventory struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*Inventory, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &Inventory{
		store: c.Store,
		index: c.Index,
	}, nil
}

func (t *Inventory) StoreContent(rc io.ReadCloser, v ct.Version, c ct.Changes) ([]string, error) {
	if rc != nil {
		defer rc.Close()

		// NOTE(leeola): a nil/zero length []byte will cause a Reader to "do nothing",
		// but not inherently return EOF (per interface spec). Therefor it should be
		// fine to read no data but test if we're at EOF, which is all this
		// ContentType can use.
		if _, err := rc.Read(nil); err != io.EOF {
			return nil, errors.New("cannot store content in the inventory type")
		}
	}

	return t.StoreMeta(v, c)
}

func (t *Inventory) StoreMeta(v ct.Version, c ct.Changes) ([]string, error) {
	var meta Meta

	if v.Meta != "" {
		if err := store.ReadAndUnmarshal(t.store, v.Meta, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// update the meta with any changes matching the meta
	meta.FromChanges(c)

	// Ensure the meta is valid. Eg, has a name, container, etc.
	if err := t.CheckMetaValid(meta); err != nil {
		return nil, errors.Stack(err)
	}

	mH, vH, err := ct.WriteMetaAndVersion(t.store, t.index, v, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return []string{mH, vH}, nil
}

func (Inventory) UnmarshalMeta(b []byte) (interface{}, error) {
	var meta Meta

	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, errors.Stack(err)
	}

	return meta, nil
}

func (Inventory) CheckMetaValid(meta Meta) error {
	// TODO(leeola): Somehow return an error that can be json printed. Ie,
	// it is intended to be shown to the api caller.
	// Likely make a UserError interface, and if it implements that we can return
	// the message, so we don't leak data.
	if meta.Name == "" {
		return errors.New("missing required inventory field: Name")
	}

	return nil
}
