package local

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"

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

func (l *Local) ReadHash(h string) (kala.Version, error) {
	var v kala.Version
	if err := ReadAndUnmarshal(l.store, h, &v); err != nil {
		return kala.Version{}, err
	}

	if structs.IsZero(v) {
		return kala.Version{}, kala.ErrNotVersion
	}

	if v.JsonHash != "" {
		if err := ReadAndUnmarshal(l.store, v.JsonHash, &v.Json); err != nil {
			return kala.Version{}, err
		}
	}

	if v.MultiBlobHash != "" {
		// TODO(leeola): Construct a new multiblob reader for the given hash.
		return kala.Version{}, errors.New("multiBlob reading not yet supported")
	}

	return v, nil
}

func (l *Local) ReadId(id string) (kala.Version, error) {
	// TODO(leeola): search the unique/id index for the given id,
	// but first i need to decide how the indexes are going to exactly
	// store the unique id versions.
	return kala.Version{}, errors.New("not implemented")
}

func (l *Local) Write(c kala.Commit, j kala.Json, r io.Reader) ([]string, error) {
	// For quicker prototyping, only supporting metadata atm
	if r != nil {
		return nil, errors.New("reader not yet implemented")
	}

	if structs.IsZero(j) && r == nil {
		return nil, errors.New("No data given to write")
	}

	jsonHash, err := MarshalAndWrite(l.store, j)
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

	versionHash, err := MarshalAndWrite(l.store, version)
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

// WriteReader writes the given reader's content to the store.
func WriteReader(s kala.Store, r io.Reader) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if r == nil {
		return "", errors.New("Reader is nil")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to readall")
	}

	h, err := s.Write(b)
	return h, errors.Wrap(err, "store failed to write")
}

// MarshalAndWrite marshals the given interface to json and writes that to the store.
func MarshalAndWrite(s kala.Store, v interface{}) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if v == nil {
		return "", errors.New("Interface is nil")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func ReadAll(s kala.Store, h string) ([]byte, error) {
	rc, err := s.Read(h)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func ReadAndUnmarshal(s kala.Store, h string, v interface{}) error {
	_, err := ReadAndUnmarshalWithBytes(s, h, v)
	return err
}

func ReadAndUnmarshalWithBytes(s kala.Store, h string, v interface{}) ([]byte, error) {
	b, err := ReadAll(s, h)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return nil, errors.Stack(err)
	}

	return b, nil
}
