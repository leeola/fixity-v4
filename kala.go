package kala

import (
	"encoding/json"
	"io"
	"time"

	"github.com/leeola/kala/q"
)

type Version struct {
	JsonHash      string `json:"metaHash,omitempty"`
	MultiBlobHash string `json:"multiBlobHash,omitempty"`

	Id                   string     `json:"id,omitempty"`
	UploadedAt           *time.Time `json:"uploadedAt,omitempty"`
	PreviousVersionCount int        `json:"previousVersionCount,omitempty"`
	PreviousVersionHash  string     `json:"previousVersion,omitempty"`

	ChangeLog string `json:"changeLog,omitempty"`

	// Json is the unmarshalled contents of the JsonHash.
	Json Json `json:"-"`

	// MultiBlob is the unmarshalled contents of the MultiBlobHash.
	//
	// This must be closed if not nil!
	MultiBlob io.ReadCloser `json:"-"`
}

type Json struct {
	Meta JsonMeta        `json:"jsonMeta,omitempty"`
	Json json.RawMessage `json:"json"`
}

type JsonMeta struct {
	IndexedFields Fields
}

type MultiBlob struct {
	BlobHashes []string `json:"blobHashes"`
}

type Blob struct {
	Blob []byte `json:"blob"`
}

type Commit struct {
	Id                  string     `json:"id,omitempty"`
	PreviousVersionHash string     `json:"previousVersion,omitempty"`
	UploadedAt          *time.Time `json:"uploadedAt,omitempty"`
	ChangeLog           string     `json:"changeLog,omitempty"`
}

// Kala implements writing, indexing and reading with a Kala store.
//
// This interface will be implemented for multiple stores, such as a local on
// disk store and a remote over network store.
type Kala interface {
	// ReadHash unmarshals the given hash contents into a Version.
	//
	// Included in the Version is the Json and MultiBlob, if any exist. If no
	// Json exists the Json struct will be zero value, and if no MultiBlob
	// exists the ReadCloser will be nil.
	//
	// ReadHash will return ErrNotVersion if the given hash is not a valid hash.
	ReadHash(hash string) (Version, error)

	// ReadId unmarshals the given id into a Version struct.
	//
	// Included in the Version is the Json and MultiBlob, if any exist. If no
	// Json exists the Json struct will be zero value, and if no MultiBlob
	// exists the ReadCloser will be nil.
	ReadId(id string) (Version, error)

	// Search for documents matching the given query.
	Search(*q.Query) ([]string, error)

	// Write the given  Commit, Meta and Reader to the Kala store.
	Write(Commit, Json, io.Reader) ([]string, error)
}
