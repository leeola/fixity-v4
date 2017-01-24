package image

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/ctutil"
	"github.com/leeola/kala/contenttype/file"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
)

type Meta struct {
	file.FileMeta

	Format string `json:"format"`
}

type Config struct {
	Store store.Store
	Index index.Indexer
}

type ContentStorer struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*ContentStorer, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &ContentStorer{
		store: c.Store,
		index: c.Index,
	}, nil
}

// TODO(leeola): centralize the common tasks in this method into helpers.
// A lot of this (writing content roller and multipart, etc) is going to be
// duplicated on every ContentType handler.
func (f *ContentStorer) StoreContent(rc io.ReadCloser, mb []byte, c contenttype.Changes) ([]string, error) {
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
	c.SetMultiPart(h)

	if err := f.index.Entry(h); err != nil {
		return nil, errors.Stack(err)
	}

	var meta Meta
	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if len(mb) == 0 {
		if h, _ := c.GetPreviousMeta(); h != "" {
			if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
				return nil, errors.Stack(err)
			}
		}
	} else {
		if err := json.Unmarshal(mb, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
	meta.FileMeta.ApplyChanges(c)
	meta.ApplyChanges(c)

	// if there is an anchor, always return the anchor for a consistent UX
	if meta.Anchor != "" {
		hashes = append(hashes, meta.Anchor)
	}

	h, err = WriteFileMeta(f.store, f.index, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)

	return hashes, nil
}

func (f *ContentStorer) StoreMeta(mb []byte, c contenttype.Changes) ([]string, error) {
	var (
		meta   Meta
		hashes []string
	)

	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if len(mb) == 0 {
		if h, _ := c.GetPreviousMeta(); h != "" {
			if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
				return nil, errors.Stack(err)
			}
		}
	} else {
		if err := json.Unmarshal(mb, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
	meta.FileMeta.ApplyChanges(c)
	meta.ApplyChanges(c)

	// if there is an anchor, always return the anchor so that the caller can easily
	// track the anchor of the content. For a consistent UX.
	if meta.Anchor != "" {
		hashes = append(hashes, meta.Anchor)
	}

	h, err := WriteFileMeta(f.store, f.index, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)

	return hashes, nil
}

func (c *ContentStorer) MetaToChanges([]byte) (contenttype.Changes, error) {
	return nil, errors.New("not implemented")
}

func WriteFileMeta(s store.Store, i index.Indexer, m Meta) (string, error) {
	// Now write the meta as well.
	h, err := store.MarshalAndWrite(s, m)
	if err != nil {
		return "", errors.Stack(err)
	}

	// Pass the changes as metadata to the indexer.
	if err := i.Metadata(h, m); err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func (m *Meta) ApplyChanges(c contenttype.Changes) {
	if f, ok := c.GetString("format"); ok {
		m.Format = f
	} else if m.Filename != "" {
		m.Format = filepath.Ext(m.Filename)

		// remove the dot from the extension
		if m.Format != "" {
			m.Format = m.Format[1:]
		}
	}
}

func (m Meta) ToMetadata(im index.Metadata) {
	if m.Format != "" {
		im["format"] = m.Format
	}
}
