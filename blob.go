package fixity

//go:generate stringer -type=BlobType -output=blob_string.go

import "io"

type BlobType int

const (
	BlobTypeSchemaless BlobType = iota
	BlobTypeParts
	BlobTypeData
	BlobTypeValues
	BlobTypeMutation
)

type BlobTyper interface {
	BlobType() (BlobType, error)
}

type BlobReadCloser interface {
	io.ReadCloser
	BlobTyper
}
