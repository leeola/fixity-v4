package store

import (
	"io"
	"time"

	"github.com/leeola/kala/index"
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

	// List records in the store
	//
	// TODO(leeola): Enable this. In the current Store implementations this is not
	// supported. However, this will be how the Indexer constructs the initial index
	// List(max, offset int) (<-chan string, error)
}

type PartRoller interface {
	Roll() (Part, error)
}

type Anchor struct {
	AnchorRand int `json:"anchorRand"`
}

type Meta struct {
	ContentType  string    `json:"contentType"`
	PreviousMeta string    `json:"previousMeta"`
	UploadedAt   time.Time `json:"uploadedAt"`

	Anchor    string `json:"anchor,omitempty"`
	MultiPart string `json:"multiPart,omitempty"`
	MultiHash string `json:"multiHash,omitempty"`

	ChangeType string `json:"changeType,omitempty"`
	ChangeLog  string `json:"changeLog,omitempty"`
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

func (m Meta) ToMetadata() index.Metadata {
	im := index.Metadata{}
	if m.Anchor != "" {
		im["anchor"] = m.Anchor
	}
	if m.MultiHash != "" {
		im["multiHash"] = m.MultiHash
	}
	if m.MultiPart != "" {
		im["multiPart"] = m.MultiPart
	}
	if !m.UploadedAt.IsZero() {
		im["uploadedAt"] = m.UploadedAt
	}
	if m.PreviousMeta != "" {
		im["previousMeta"] = m.PreviousMeta
	}
	return im
}
