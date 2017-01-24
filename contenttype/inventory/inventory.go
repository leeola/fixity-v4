package inventory

import (
	"encoding/json"
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
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

// TODO(leeola): centralize the common tasks in this method into helpers.
// A lot of this (writing content roller and multipart, etc) is going to be
// duplicated on every ContentType handler.
func (i *Inventory) StoreContent(rc io.ReadCloser, mb []byte, c contenttype.Changes) ([]string, error) {
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

	return i.StoreMeta(mb, c)
}

func (i *Inventory) StoreMeta(mb []byte, c contenttype.Changes) ([]string, error) {
	var (
		meta   Meta
		hashes []string
	)

	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if len(mb) == 0 {
		if h, _ := c.GetPreviousMeta(); h != "" {
			if err := store.ReadAndUnmarshal(i.store, h, &meta); err != nil {
				return nil, errors.Stack(err)
			}
		}
	} else {
		if err := json.Unmarshal(mb, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// Apply the filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
	meta.FromChanges(c)

	// Ensure the meta is valid. Eg, has a name, container, etc.
	if err := i.CheckMetaValid(meta); err != nil {
		return nil, errors.Stack(err)
	}

	// if there is an anchor, always return the anchor so that the caller can easily
	// track the anchor of the content. For a consistent UX.
	if meta.Anchor != "" {
		hashes = append(hashes, meta.Anchor)
	}

	h, err := contenttype.WriteMeta(i.store, i.index, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)

	return hashes, nil
}

func (i *Inventory) MetaToChanges(b []byte) (contenttype.Changes, error) {
	m := Meta{}
	if len(b) != 0 {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, errors.Stack(err)
		}
	}

	c := contenttype.Changes{}
	m.ToChanges(c)
	return c, nil
}

func (i *Inventory) CheckMetaValid(meta Meta) error {
	// TODO(leeola): Somehow return an error that can be json printed. Ie,
	// it is intended to be shown to the api caller.
	// Likely make a UserError interface, and if it implements that we can return
	// the message, so we don't leak data.
	if meta.Name == "" {
		return errors.New("missing required inventory field: Name")
	}

	if meta.Container != "" {
		valid, err := store.IsAnchor(i.store, meta.Container)
		if err != nil {
			return errors.Wrap(err, "failed to read Store to check if Container is an Anchor")
		}
		if !valid {
			return errors.New("Container value must be the hash of an existing anchor")
		}
	}

	return nil
}
