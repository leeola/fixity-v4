package store

import "io"

type Store interface {
	// Check if the given hash exists in the Store
	Exists(string) (bool, error)

	// Takes a hex string of the content hash, and returns a reader for the content
	Read(string) (io.ReadCloser, error)

	// Write raw data to the store.
	//
	// Return the hash of the written data.
	Write([]byte) (string, error)

	// Write the given data to the store only if it matches the given hash.
	//
	// Note that this must compute the hash to ensure the bytes match the given hex
	// hash.
	WriteHash(string, []byte) error

	// List records in the store
	//
	// TODO(leeola): Enable this. In the current Store implementations this is not
	// supported. However, this will be how the Indexer constructs the initial index
	// List(max, offset int) (<-chan string, error)
}

type ContentRoller interface {
	Roll() (Content, error)
}

type Perma struct {
	PermaRand int `json:"permaRand"`
}

// MultiPart is a series of hashes for a single piece of data.
type MultiPart struct {
	Perma string   `json:"perma,omitempty"`
	Parts []string `json:"parts"`
}

type Content struct {
	Content []byte `json:"content"`
}
