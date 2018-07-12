package fixity

//go:generate stringer -type=BlobType -output=blob_string.go

import (
	"context"
	"fmt"
	"io"

	"github.com/leeola/fixity/config"
)

type BlobType int

const (
	BlobTypeSchemaless BlobType = iota
	BlobTypeParts
	BlobTypeData
	BlobTypeValues
	BlobTypeMutation
)

type Blobstore interface {
	BlobReader
	BlobWriter
}

type BlobTyper interface {
	BlobType() (BlobType, error)
}

type BlobReadCloser interface {
	io.ReadCloser
	BlobTyper
}

type BlobWriter interface {
	Write(context.Context, []byte) (Ref, error)
}

type BlobReader interface {
	Read(context.Context, Ref) (io.ReadCloser, error)
}

func NewBlobstoreFromConfig(name string, c config.Config) (Blobstore, error) {
	if name == "" {
		return nil, fmt.Errorf("empty blobstore name")
	}

	tc, ok := c.BlobstoreConfigs[name]
	if !ok {
		return nil, fmt.Errorf("blobstore name not found: %q", name)
	}

	constructor, ok := blobstoreRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("blobstore type not found: %q", tc.Type)
	}

	bs, err := constructor.New(name, c)
	if err != nil {
		return nil, fmt.Errorf("blobstore constructor %s: %v", name, err)
	}

	return bs, nil
}
