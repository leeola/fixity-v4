package file

import (
	"encoding/json"
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/ctutil"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
)

type FileMeta struct {
	contenttype.Meta

	Filename string `json:"filename"`
}

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

// TODO(leeola): centralize the common tasks in this method into helpers.
// A lot of this (writing content roller and multipart, etc) is going to be
// duplicated on every ContentType handler.
func (f *File) StoreContent(rc io.ReadCloser, mb []byte, c contenttype.Changes) ([]string, error) {
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

	// Write the entries, not including the final metadata hash
	// The last hash is metadata, and we'll add that manually.
	for _, h := range hashes {
		if err := f.index.Entry(h); err != nil {
			return nil, errors.Stack(err)
		}
	}

	var meta FileMeta
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

	// Apply the filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
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

func (f *File) StoreMeta(mb []byte, c contenttype.Changes) ([]string, error) {
	var (
		meta   FileMeta
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

	// Apply the filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
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

func (f *File) MetaToChanges(b []byte) (contenttype.Changes, error) {
	m := FileMeta{}
	if len(b) != 0 {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, errors.Stack(err)
		}
	}

	c := contenttype.Changes{}
	m.ToChanges(c)
	return c, nil
}

func WriteFileMeta(s store.Store, i index.Indexer, m FileMeta) (string, error) {
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

func (m *FileMeta) ApplyChanges(c contenttype.Changes) {
	m.Meta.FromChanges(c)

	if f, ok := c.GetString("filename"); ok {
		m.Filename = f
	}
	// If the name wasn't explicitly set (ie, to an empty value), and
	// the name is empty and the filename is not empty, default the name to
	// the filename.
	if _, ok := c.GetString("name"); !ok && m.Name == "" && m.Filename != "" {
		m.Name = m.Filename
	}
}

func (m FileMeta) ToMetadata(im index.Metadata) {
	if m.Filename != "" {
		im["filename"] = m.Filename
	}
}

func (m FileMeta) ToChanges(c contenttype.Changes) {
	m.Meta.ToChanges(c)
	if m.Filename != "" {
		c.Set("filename", m.Filename)
	}
}

func UnmarshalMetadata(b []byte) (index.Indexable, error) {
	m := FileMeta{}

	if len(b) != 0 {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, errors.Stack(err)
		}
	}

	return m, nil
}
