package fixity

import (
	"io"
)

// WriteRequest represents a blob to be written alone with metadata.
type WriteRequest struct {
	Id       string
	RollSize int
	Fields   Fields
	Blob     io.ReadCloser
}

// NewWrite creates a new WriteRequest to be used with Fixity.WriteRequest.
func NewWrite(id string, rc io.ReadCloser, f ...Field) *WriteRequest {
	return &WriteRequest{
		Id:     id,
		Fields: f,
		Blob:   rc,
	}
}
