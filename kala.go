package kala

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type Commit struct {
	Id                  string    `json:"id,omitempty"`
	PreviousVersionHash string    `json:"previousVersion,omitempty"`
	UploadedAt          time.Time `json:"uploadedAt,omitempty"`
	ChangeLog           string    `json:"changeLog,omitempty"`
}

// Meta is a type alias for store.Meta for the UX of Kala package users.
type Meta store.Meta

type Config struct {
	Index index.Index
	Store store.Store
}

type Kala struct {
	config Config
	index  index.Index
	store  store.Store
}

func New(c Config) (*Kala, error) {
	if c.Index == nil {
		return nil, errors.New("missing reqired config: Index")
	}
	if c.Store == nil {
		return nil, errors.New("missing reqired config: Store")
	}

	return &Kala{
		config: c,
		index:  c.Index,
		store:  c.Store,
	}, nil
}

func (k *Kala) Write(c Commit, m Meta, r io.Reader) ([]string, error) {
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
