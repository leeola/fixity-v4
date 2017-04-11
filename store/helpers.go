package store

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
)

// WriteReader writes the given reader's content to the store.
func WriteReader(s Store, r io.Reader) (string, error) {
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
func MarshalAndWrite(s Store, v interface{}) (string, error) {
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

// WriteBlob is a type safe version of MarshalAndWrite.
func WriteBlob(s Store, v Blob) (string, error) {
	return MarshalAndWrite(s, v)
}

// WriteMeta is a type safe version of MarshalAndWrite.
func WriteMeta(s Store, v Meta) (string, error) {
	return MarshalAndWrite(s, v)
}

// WriteMultiBlob is a type safe version of MarshalAndWrite.
func WriteMultiBlob(s Store, v Meta) (string, error) {
	return MarshalAndWrite(s, v)
}

// WriteVersion is a type safe version of MarshalAndWrite.
func WriteVersion(s Store, v Version) (string, error) {
	return MarshalAndWrite(s, v)
}
