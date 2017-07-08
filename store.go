package fixity

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
)

type Store interface {
	// Check if the given hash exists in the Store
	Exists(string) (bool, error)

	// Takes a hex string of the content hash, and returns a reader for the content
	Read(string) (io.ReadCloser, error)

	// Write raw data to the store.
	//
	// The created return value allows the caller to know if the written content
	// was already in the store or not. Pre-existing content can signal whether
	// or not a Fixity block is to be created or not.
	//
	// Return the hash of the written data.
	Write([]byte) (hash string, created bool, err error)

	// Write the given data to the store only if it matches the given hash.
	//
	// This is used for Node->Node or Store->Store writes, and ensures a write
	// is only ever done if the hash matches.
	//
	// Like Write(), WriteHash() returns a bool indicating whether the write was
	// created or not.
	//
	// Note that this must compute the hash to ensure the bytes match the given
	// hash.
	WriteHash(string, []byte) (created bool, err error)

	// List records in the store.
	//
	// IMPORTANT: Listing may not be deterministic and does not ensure that new records
	// or removed records are included in the listing. Therefor Listing should be done
	// before before a store is being actively served.
	List() (<-chan string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

// Chunker implements chunking over bytes.
type Chunker interface {
	Chunk() (Chunk, error)
}

func writeReader(s Store, r io.Reader) (string, error) {
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

	h, _, err := s.Write(b)
	return h, errors.Wrap(err, "store failed to write")
}

func marshalAndWrite(s Store, v interface{}) (string, error) {
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

	h, _, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func readAll(s Store, h string) ([]byte, error) {
	rc, err := s.Read(h)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func readAndUnmarshal(s Store, h string, v interface{}) error {
	_, err := readAndUnmarshalWithBytes(s, h, v)
	return err
}

func readAndUnmarshalWithBytes(s Store, h string, v interface{}) ([]byte, error) {
	b, err := readAll(s, h)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return nil, errors.Stack(err)
	}

	return b, nil
}
