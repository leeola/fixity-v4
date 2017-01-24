package handlers

import "github.com/leeola/kala/contenttype"

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
	Changes contenttype.Changes `json:"changes"`
	Error   string              `json:"error,omitempty"`
}
