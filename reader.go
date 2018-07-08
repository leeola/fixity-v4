package fixity

import "io"

type Reader interface {
	io.Reader

	Size() (int64, error)

	Checksum() (string, error)
}
