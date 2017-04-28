package local

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/fatih/structs"
	"github.com/leeola/errors"
	"github.com/leeola/kala"
	"github.com/leeola/kala/q"
)

type Config struct {
	Index kala.Index
	Store kala.Store
}

type Local struct {
	config Config
	index  kala.Index
	store  kala.Store
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

func (l *Local) Write(c kala.Commit, j kala.Json, r io.Reader) ([]string, error) {
	// For quicker prototyping, only supporting metadata atm
	if r != nil {
		return nil, errors.New("reader not yet implemented")
	}

	if structs.IsZero(j) && r == nil {
		return nil, errors.New("No data given to write")
	}

	jsonHash, err := kala.MarshalAndWrite(l.store, j)
	if err != nil {
		return nil, errors.Stack(err)
	}

	var multiBlobHash string
	// TODO(leeola): Make this into a multipart splitter.
	// For now it's disabled.
	//
	// multiBlobHash, err := store.WriteReader(l.store, r)
	// if err != nil {
	// return nil, errors.Stack(err)
	// }

	version := kala.Version{
		JsonHash:      jsonHash,
		MultiBlobHash: multiBlobHash,
	}

	// TODO(leeola): load the old version if previous version hash is specified
	// if c.PreviousVersionHash != "" {
	// // .. load previous hash
	// version = previousVersion
	// }

	versionHash, err := kala.MarshalAndWrite(l.store, version)
	if err != nil {
		return nil, errors.Stack(err)
	}

	// TODO(leeola): Index the metadata now that all has been written to the store.

	// Replace the old changelog no matter what. Eg, even if we loaded an old version,
	// the old version's changelog doesn't apply to the new version, so replace it,
	// even if we're repalcing it with nothing.
	version.ChangeLog = c.ChangeLog

	var hashes []string
	if jsonHash != "" {
		hashes = append(hashes, jsonHash)
	}

	// copy the fields list so that we can add to it, without
	// modifying what is stored
	indexFields := make(kala.Fields, len(j.Meta.IndexedFields))
	for i, f := range j.Meta.IndexedFields {
		indexFields[i] = f

		// TODO(leeola): Check for nil field values, and attempt to find
		// the value by unmarshalling the json manually.
		//
		// This should be done in this location, and not modifying the actual
		// j.Meta.IndexedFields or we would end up storing values twice.
		if f.Value == nil {
			return nil, errors.New("automatic json field value assertion not yet supported")
		}
	}

	indexFields.Append(kala.Field{
		Field: "version.jsonHash",
		Value: version.JsonHash,
	})
	indexFields.Append(kala.Field{
		Field: "version.multiBlobHash",
		Value: version.MultiBlobHash,
	})
	indexFields.Append(kala.Field{
		Field: "version.id",
		Value: version.Id,
	})
	indexFields.Append(kala.Field{
		Field: "version.uploadedAt",
		Value: version.UploadedAt,
	})
	indexFields.Append(kala.Field{
		Field: "version.previousVersionCount",
		Value: version.PreviousVersionCount,
	})
	indexFields.Append(kala.Field{
		Field: "version.previousVersionHash",
		Value: version.PreviousVersionHash,
	})

	if err := l.index.Index(versionHash, version.Id, indexFields); err != nil {
		return nil, err
	}

	return append(hashes, versionHash), nil
}

func (l *Local) Search(q *q.Query) ([]string, error) {
	return l.index.Search(q)
}

// NewId is a helper to generate a new default length Id.
//
// Note that the Id is encoded as hex to easily interact with it, rather
// than plain bytes.
func NewId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
