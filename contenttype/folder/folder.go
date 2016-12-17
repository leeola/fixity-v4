package folder

import (
	"encoding/json"
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type FolderMeta struct {
	store.Meta

	Foldername string `json:"foldername"`
}

type Config struct {
	Store store.Store
	Index index.Indexer
}

type Folder struct {
	store store.Store
	index index.Indexer
}

func New(c Config) (*Folder, error) {
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}
	if c.Index == nil {
		return nil, errors.New("missing required Config field: Index")
	}

	return &Folder{
		store: c.Store,
		index: c.Index,
	}, nil
}

func (f *Folder) StoreContent(rc io.ReadCloser, mb []byte, c contenttype.Changes) ([]string, error) {
	// Folder doesn't allow content, so close any reader given and error.
	if rc != nil {
		rc.Close()
		return nil, errors.New("Folder ContentType cannot store MultiPart data")
	}

	return f.Meta(mb, c)
}

func (f *Folder) Meta(mb []byte, c contenttype.Changes) ([]string, error) {
	var (
		meta   FolderMeta
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
	// This maps the fields in the MetaChanges map to the Meta and FileMeta struct.
	contenttype.ApplyCommonChanges(&meta.Meta, c)
	meta.ApplyChanges(c)

	var multiHash store.MultiHash
	if multiHashes, ok := c["hash"]; ok {
		// the caller specified the hashes to use manually, so overwrite the loaded
		// hash values.
		multiHash.Hashes = multiHashes
	} else if mh := meta.MultiHash; mh != "" {
		// the caller didn't explicitly set the hashes, so load the previous hashes
		// if a multihash was specified.
		if err := store.ReadAndUnmarshal(f.store, mh, &multiHash); err != nil {
			return nil, errors.Stack(err)
		}
	}

	// add any hashes specified
	if ah, _ := c["addHash"]; len(ah) > 0 {
		multiHash.Hashes = append(multiHash.Hashes, ah...)
	}

	if dhs, ok := c["delHash"]; ok {
		for _, dh := range dhs {
			for i, h := range multiHash.Hashes {
				if dh == h {
					multiHash.Hashes = append(multiHash.Hashes[:i], multiHash.Hashes[i+1:]...)
				}
			}
		}
	}

	h, err := store.WriteMultiHash(f.store, multiHash)
	if err != nil {
		return nil, errors.Stack(err)
	}
	meta.MultiHash = h
	hashes = append(hashes, h)

	// if there is an anchor, always return the anchor so that the caller can easily
	// track the anchor of the content. For a consistent UX.
	if meta.Anchor != "" {
		hashes = append(hashes, meta.Anchor)
	}

	h, err = WriteMeta(f.store, f.index, meta)
	if err != nil {
		return nil, errors.Stack(err)
	}
	hashes = append(hashes, h)

	return hashes, nil
}

func WriteMeta(s store.Store, i index.Indexer, m FolderMeta) (string, error) {
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

func (m *FolderMeta) ApplyChanges(c contenttype.Changes) {
	if n, ok := c.GetString("foldername"); ok {
		m.Foldername = n
	}
}

func (m FolderMeta) ToMetadata() index.Metadata {
	im := m.Meta.ToMetadata()
	if m.Foldername != "" {
		im["foldername"] = m.Foldername
	}
	return im
}
