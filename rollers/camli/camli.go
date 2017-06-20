package camli

import (
	"bufio"
	"bytes"
	"io"

	"camlistore.org/pkg/rollsum"

	"github.com/leeola/errors"
	"github.com/leeola/fixity"
)

type Roller struct {
	reader      *bufio.Reader
	rollSum     *rollsum.RollSum
	minRollSize int64
}

func New(r io.Reader, rollSize int64) (*Roller, error) {
	if r == nil {
		return nil, errors.New("missing Reader")
	}

	return &Roller{
		// TODO(leeola): explicitly declare the buffer size.
		reader:      bufio.NewReader(r),
		rollSum:     rollsum.New(),
		minRollSize: rollSize,
	}, nil
}

func (r *Roller) Roll() (fixity.Chunk, error) {
	var (
		byteContent bytes.Buffer
		byteCount   int64
	)

	// TODO(leeola): Add a peek method to break out of the loop if the end of the
	// roller is near. This way we don't create small tailing chunks if possible.
	for {
		if r.reader == nil {
			return fixity.Chunk{}, io.EOF
		}

		c, err := r.reader.ReadByte()
		if err != nil && err != io.EOF {
			return fixity.Chunk{}, errors.Stack(err)
		}

		// if we're EOF, break so we can return the existing content.
		if err == io.EOF {
			// nil the reader so that on the next Roll, we'll return EOF
			r.reader = nil
			break
		}

		byteCount = byteCount + 1
		if err := byteContent.WriteByte(c); err != nil {
			return fixity.Chunk{}, errors.Stack(err)
		}

		r.rollSum.Roll(c)

		if r.rollSum.OnSplit() && byteCount > r.minRollSize {
			break
		}
	}

	return fixity.Chunk{
		ChunkBytes: byteContent.Bytes(),
		Size:       byteCount,
	}, nil
}
