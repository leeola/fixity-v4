package handlers

import (
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/store"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type HashResponse struct {
	Hash  string `json:"hash,omitempty"`
	Error string `json:"error,omitempty"`
}

type HashesResponse struct {
	Hashes []string `json:"hashes,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type BlobContentTypeResponse struct {
	ContentType string `json:contentType,omitempty`
	Error       string `json:"error,omitempty"`
}

type ChangesResponse struct {
	Changes ct.Changes `json:"changes"`
	Error   string     `json:"error,omitempty"`
}

type ResolveBlobResponse struct {
	Hash  string `json:"hash"`
	Blob  []byte `json:"blob"`
	Error string `json:"error,omitempty"`
}

type ResolveVersionResponse struct {
	Hash    string        `json:"hash"`
	Version store.Version `json:"version"`
	Error   string        `json:"error,omitempty"`
}
