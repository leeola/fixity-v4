package store

import (
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

type PartRoller interface {
	Roll() (Part, error)
}

type Version struct {
	ContentType string `json:"contentType"`
	Meta        string `json:"meta"`

	Anchor               string    `json:"anchor,omitempty"`
	UploadedAt           time.Time `json:"uploadedAt,omitempty"`
	PreviousVersionCount int       `json:"previousVersionCount"`
	PreviousVersion      string    `json:"previousVersion"`

	ChangeLog string `json:"changeLog,omitempty"`
}

type Meta struct {
	MultiPart string `json:"multiPart,omitempty"`
	MultiHash string `json:"multiHash,omitempty"`

	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type MultiHash struct {
	Hashes []string `json:"hashes"`
}

type MultiPart struct {
	Parts []string `json:"parts"`
}

type Part struct {
	Part []byte `json:"part"`
}
