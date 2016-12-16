package file

import (
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
)

type FileMeta struct {
	store.Meta

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
func (f *File) StoreContent(rc io.ReadCloser, c store.MetaChanges) ([]string, error) {
	if rc == nil {
		return nil, errors.New("missing ReadCloser")
	}
	defer rc.Close()

	roller, err := camli.New(rc)
	if err != nil {
		return nil, errors.Stack(err)
	}

	// write the actual content
	hashes, err := store.WriteContentRoller(f.store, roller)
	if err != nil {
		return nil, errors.Stack(err)
	}

	var meta FileMeta

	// write the multipart
	h, err := store.WriteMultiPart(f.store, store.MultiPart{
		Parts: hashes,
	})
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)
	c.SetMulti(h)

	// Write the entries, not including the final metadata hash
	// The last hash is metadata, and we'll add that manually.
	for _, h := range hashes {
		if err := f.index.Entry(h); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// If the previous hash exists, populate the above filemeta with the data
	// in the hash.
	if h, _ := c.GetPreviousMeta(); h != "" {
		if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the MetaChanges map to the Meta and FileMeta struct.
	store.ApplyCommonChanges(&meta.Meta, c)
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

func (f *File) Meta(c store.MetaChanges) ([]string, error) {
	var (
		meta   FileMeta
		hashes []string
	)

	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if h, _ := c.GetPreviousMeta(); h != "" {
		if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the MetaChanges map to the Meta and FileMeta struct.
	store.ApplyCommonChanges(&meta.Meta, c)
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

func (m *FileMeta) ApplyChanges(c store.MetaChanges) {
	if f, ok := c.GetString("filename"); ok {
		m.Filename = f
	}
}

func (m FileMeta) ToMetadata() index.Metadata {
	im := m.Meta.ToMetadata()
	if m.Filename != "" {
		im["filename"] = m.Filename
	}
	if m.ContentType != "" {
		im["contentType"] = m.ContentType
	}
	return im
}
