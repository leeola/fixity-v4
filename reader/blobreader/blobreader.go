package blobreader

import (
	"fmt"
	"io"

	"github.com/leeola/fixity"
)

// ReadCloser implements peek/buffer into a BlobType to identify the blob,
// and then caches the result for repeated BlobType requests.
type ReadCloser struct {
	io.ReadCloser
}

func New(rc io.ReadCloser) *ReadCloser {
	return &ReadCloser{ReadCloser: rc}
}

func (rc *ReadCloser) BlobType() (fixity.BlobType, error) {
	return 0, fmt.Errorf("not implemented")
}
