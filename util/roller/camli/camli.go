package camli

import (
	"bufio"
	"bytes"
	"io"

	"camlistore.org/pkg/rollsum"

	"github.com/leeola/errors"
	"github.com/leeola/fixity/store"
)

const (
	// The arbitrarily chosen default chunk size. Ideally it should strike a balance
	// between too many chunks written/indexed, and still allowing frequent changes
	// to files to reduce chunk size.
	//
	// If a file format is commonly editing the first X bytes, this value should likely
	// be overridden in the Roller itself.
	MinRollSize = 4096000
)

type Roller struct {
	reader      *bufio.Reader
	rollSum     *rollsum.RollSum
	MinRollSize int
}

func New(r io.Reader) (*Roller, error) {
	if r == nil {
		return nil, errors.New("missing Reader")
	}

	return &Roller{
		// TODO(leeola): explicitly declare the buffer size.
		reader:      bufio.NewReader(r),
		rollSum:     rollsum.New(),
		MinRollSize: MinRollSize,
	}, nil
}

func (r *Roller) Roll() (store.Part, error) {
	var (
		byteContent bytes.Buffer
		byteCount   int
	)

	// TODO(leeola): Add a peek method to break out of the loop if the end of the
	// roller is near. This way we don't create small tailing chunks if possible.
	for {
		if r.reader == nil {
			return store.Part{}, io.EOF
		}

		c, err := r.reader.ReadByte()
		if err != nil && err != io.EOF {
			return store.Part{}, errors.Stack(err)
		}

		// if we're EOF, break so we can return the existing content.
		if err == io.EOF {
			// nil the reader so that on the next Roll, we'll return EOF
			r.reader = nil
			break
		}

		byteCount = byteCount + 1
		if err := byteContent.WriteByte(c); err != nil {
			return store.Part{}, errors.Stack(err)
		}

		r.rollSum.Roll(c)

		if r.rollSum.OnSplit() && byteCount > r.MinRollSize {
			break
		}
	}

	return store.Part{Part: byteContent.Bytes()}, nil
}
