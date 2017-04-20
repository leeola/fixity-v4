package local

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/asdine/storm/index"
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

type Config struct {
	Index index.Index
	Store store.Store
}

type Local struct {
	config Config
	index  index.Index
	store  store.Store
}

func New(c Config) (*Local, error) {
	if c.Index == nil {
		return nil, errors.New("missing reqired config: Index")
	}
	if c.Store == nil {
		return nil, errors.New("missing reqired config: Store")
	}

	return &Local{
		config: c,
		index:  c.Index,
		store:  c.Store,
	}, nil
}

func (k *Local) Write(c Commit, m Meta, r io.Reader) ([]string, error) {
	// For quicker prototyping, only supporting metadata atm
	if r != nil {
		return nil, errors.New("reader not yet implemented")
	}

	if structs.IsZero(m) && r == nil {
		return nil, errors.New("No data given to write")
	}

	metaHash, err := store.MarshalAndWrite(k.store, store.Meta(m))
	if err != nil {
		return nil, errors.Stack(err)
	}

	var multiBlobHash string
	// TODO(leeola): Make this into a multipart splitter
	// multiBlobHash, err := store.WriteReader(k.store, r)
	// if err != nil {
	// return nil, errors.Stack(err)
	// }

	version := store.Version{
		MetaHash:      metaHash,
		MultiBlobHash: multiBlobHash,
	}

	// TODO(leeola): load the old version if previous version hash is specified
	// if c.PreviousVersionHash != "" {
	// // .. load previous hash
	// version = previousVersion
	// }

	versionHash, err := store.WriteVersion(s, version)
	if err != nil {
		return nil, errors.Stack(err)
	}

	// TODO(leeola): Index the metadata now that all has been written to the store.

	// Replace the old changelog no matter what. Eg, even if we loaded an old version,
	// the old version's changelog doesn't apply to the new version, so replace it,
	// even if we're repalcing it with nothing.
	version.ChangeLog = c.ChangeLog

	var hashes []string
	if metaHash != "" {
		hashes = append(hashes, metaHash)
	}

	return append(hashes, versionHash), nil
}

// NewId is a helper to generate a new default length Id.
//
// Note that the Id is encoded as hex to easily interact with it, rather
// than plain bytes.
func NewId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
