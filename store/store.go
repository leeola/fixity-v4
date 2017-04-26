package store

import (
	"encoding/json"
	"io"
	"time"
)

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

	// List records in the store.
	//
	// IMPORTANT: Listing may not be deterministic and does not ensure that new records
	// or removed records are included in the listing. Therefor Listing should be done
	// before before a store is being actively served.
	List() (<-chan string, error)
}

type Version struct {
	JsonHash      string `json:"metaHash,omitempty"`
	MultiBlobHash string `json:"multiBlobHash,omitempty"`

	Id                   string    `json:"id,omitempty"`
	UploadedAt           time.Time `json:"uploadedAt,omitempty"`
	PreviousVersionCount int       `json:"previousVersionCount,omitempty"`
	PreviousVersionHash  string    `json:"previousVersion,omitempty"`

	ChangeLog string `json:"changeLog,omitempty"`
}

type Json struct {
	Meta JsonMeta        `json:"jsonMeta,omitempty"`
	Json json.RawMessage `json:"json"`
}

type JsonMeta struct {
	IndexedFields []Field
}

type MultiBlob struct {
	BlobHashes []string `json:"blobHashes"`
}

type Blob struct {
	Blob []byte `json:"blob"`
}
