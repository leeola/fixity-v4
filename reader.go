package fixity

import "io"

type Reader interface {
	io.Reader

	Size() int64

	Checksum() string
}
