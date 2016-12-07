package camli

import (
	"bufio"
	"bytes"
	"io"

	"camlistore.org/pkg/rollsum"

	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

type Roller struct {
	reader  *bufio.Reader
	rollSum *rollsum.RollSum
}

func New(r io.Reader) (*Roller, error) {
	if r == nil {
		return nil, errors.New("missing Reader")
	}

	return &Roller{
		// TODO(leeola): explicitly declare the buffer size.
		reader:  bufio.NewReader(r),
		rollSum: rollsum.New(),
	}, nil
}

func (r *Roller) Roll() (store.Content, error) {
	var byteContent bytes.Buffer
	for {
		if r.reader == nil {
			return store.Content{}, io.EOF
		}

		c, err := r.reader.ReadByte()
		if err != nil && err != io.EOF {
			return store.Content{}, errors.Stack(err)
		}

		// if we're EOF, break so we can return the existing content.
		if err == io.EOF {
			// nil the reader so that on the next Roll, we'll return EOF
			r.reader = nil
			break
		}

		if err := byteContent.WriteByte(c); err != nil {
			return store.Content{}, errors.Stack(err)
		}

		r.rollSum.Roll(c)

		if r.rollSum.OnSplit() {
			break
		}
	}

	return store.Content{Content: byteContent.Bytes()}, nil
}
