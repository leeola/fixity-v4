package store

import "io"

const (
	PermaType     = "perma"
	ContentType   = "content"
	MultiPartType = "multiPart"
)

type Store interface {
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
	List(max, offset int) (<-chan string, error)
}

// type ContentType struct {
// 	Type string `json:"type"`
// }

type Perma struct {
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt,omitempty"`
	Rand      []byte `json:"rand"`
}

// MultiPart is a series of hashes for a single piece of data.
type MultiPart struct {
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt,omitempty"`
	Perma     string `json:"perma,omitempty"`

	// A sum of the *content* of all of the parts, combined.
	//
	// This allows for easy referencing of the data via the original checksum
	PartsSum string   `json:"PartsSum"`
	Parts    []string `json:"parts"`
}

type Content struct {
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt,omitempty"`
	Perma     string `json:"perma,omitempty"`
	Content   []byte `json:"content"`
}
