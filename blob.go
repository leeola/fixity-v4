package fixity

import "io"

type BlobType int

const (
	BlobTypeMutation = iota + 1
	BlobTypeContent
	BlobTypeParts
	BlobTypePartBytes
)

type BlobTyper interface {
	BlobType() (BlobType, error)
}

type BlobReadCloser interface {
	io.ReadCloser
	BlobTyper
}
