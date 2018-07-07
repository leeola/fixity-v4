package fixity

import "io"

type ReadCloser interface {
	io.ReadCloser

	Size() int64

	Checksum() string
}
