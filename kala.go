package kala

import (
	"io"
	"time"
)

// Kala implements writing, indexing and reading with a Kala store.
//
// This interface will be implemented for multiple stores, such as a local on
// disk store and a remote over network store.
type Kala interface {
	// Write the given  Commit, Meta and Reader to the Kala store.
	Write(Commit, Json, io.Reader) ([]string, error)
}

type Commit struct {
	Id                  string    `json:"id,omitempty"`
	PreviousVersionHash string    `json:"previousVersion,omitempty"`
	UploadedAt          time.Time `json:"uploadedAt,omitempty"`
	ChangeLog           string    `json:"changeLog,omitempty"`
}
