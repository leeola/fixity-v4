package blobstore

import (
	"context"
	"io"

	"github.com/leeola/fixity"
)

type Reader interface {
	Read(context.Context, fixity.Ref) (io.ReadCloser, error)
}

type Writer interface {
	Write(context.Context, []byte) (fixity.Ref, error)
}

type ReadWriter interface {
	Reader
	Writer
}
