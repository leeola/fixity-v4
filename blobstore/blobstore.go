package blobstore

import (
	"context"

	"github.com/leeola/fixity"
)

type Writer interface {
	Write(context.Context, []byte) (fixity.Ref, error)
}

type ReadWriter interface {
	Reader
	Writer
}
