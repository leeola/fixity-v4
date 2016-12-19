package file

import (
	"encoding/json"
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	ct "github.com/leeola/kala/contenttype"
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

func (f *File) StoreContent(rc io.ReadCloser, mb []byte, c ct.Changes) <-chan ct.Result {
	results := make(chan ct.Result, 1)
	go func() {
		if err := f.storeContent(rc, mb, c, results); err != nil {
			results <- ct.Result{Error: err}
			close(results)
		}
	}()
	return results
}

// TODO(leeola): centralize the common tasks in this method into helpers.
// A lot of this (writing content roller and multipart, etc) is going to be
// duplicated on every ContentType handler.
func (f *File) storeContent(rc io.ReadCloser, mb []byte, c contenttype.Changes, ch chan ct.Result) error {
	if rc == nil {
		return errors.New("missing ReadCloser")
	}
	defer rc.Close()

	roller, err := camli.New(rc)
	if err != nil {
		return errors.Stack(err)
	}

	// write the actual content
	parts, err := ct.WritePartRoller(f.store, f.index, roller, ch)
	if err != nil {
		return errors.Stack(err)
	}

	// write the multipart
	h, err := store.WriteMultiPart(f.store, store.MultiPart{
		Parts: parts,
	})
	if err != nil {
		return errors.Stack(err)
	}
	if err := f.index.Entry(h); err != nil {
		return errors.Stack(err)
	}
	ch <- ct.Result{Hash: h}
	c.SetMultiPart(h)

	var meta FileMeta
	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if len(mb) == 0 {
		if h, _ := c.GetPreviousMeta(); h != "" {
			if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
				return errors.Stack(err)
			}
		}
	} else {
		if err := json.Unmarshal(mb, &meta); err != nil {
			return errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
	contenttype.ApplyCommonChanges(&meta.Meta, c)
	meta.ApplyChanges(c)

	// if there is an anchor, always return the anchor for a consistent UX
	if meta.Anchor != "" {
		ch <- ct.Result{Hash: meta.Anchor}
	}

	h, err = WriteFileMeta(f.store, f.index, meta)
	if err != nil {
		return errors.Stack(err)
	}
	ch <- ct.Result{Hash: meta.Anchor}

	close(ch)
	return nil
}

func (f *File) Meta(mb []byte, c contenttype.Changes) <-chan contenttype.Result {
	results := make(chan contenttype.Result, 1)
	go func() {
		if err := f.meta(mb, c, results); err != nil {
			results <- ct.Result{Error: err}
		}
		close(results)
	}()
	return results
}

func (f *File) meta(mb []byte, c contenttype.Changes, ch chan contenttype.Result) error {
	var meta FileMeta

	// If the previous hash exists, load that metadata hash and populate the above
	// filemeta with the data in the hash.
	if len(mb) == 0 {
		if h, _ := c.GetPreviousMeta(); h != "" {
			if err := store.ReadAndUnmarshal(f.store, h, &meta); err != nil {
				return errors.Stack(err)
			}
		}
	} else {
		if err := json.Unmarshal(mb, &meta); err != nil {
			return errors.Stack(err)
		}
	}

	// Apply the common and filemeta changes to the metadata.
	// This maps the fields in the Changes map to the Meta and FileMeta struct.
	contenttype.ApplyCommonChanges(&meta.Meta, c)
	meta.ApplyChanges(c)

	// if there is an anchor, always return the anchor so that the caller can easily
	// track the anchor of the content. For a consistent UX.
	if meta.Anchor != "" {
		ch <- ct.Result{Hash: meta.Anchor}
	}

	h, err := WriteFileMeta(f.store, f.index, meta)
	if err != nil {
		return errors.Stack(err)
	}
	ch <- ct.Result{Hash: h}

	return nil
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
	if f, ok := c.GetString("filename"); ok {
		m.Filename = f
	}
}

func (m FileMeta) ToMetadata() index.Metadata {
	im := m.Meta.ToMetadata()
	if m.Filename != "" {
		im["filename"] = m.Filename
	}
	return im
}
