package restic

import (
	"io"
	"math"

	"github.com/leeola/chunker"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
)

// averageBits of 14 is resulting in roughly chunks of 12KiB by default.
//
// This size should be low enough to allow small writes to be split up,
// even at the cost of write performance (within reason).
const averageBits = 14

const minSizeMul = 0.5
const maxSizeMul = 1.5

type Roller struct {
	buf     []byte
	chunker *chunker.Chunker
}

func New(r io.Reader, averageChunkSize uint64) (*Roller, error) {
	if r == nil {
		return nil, errors.New("missing Reader")
	}

	return &Roller{
		// does this size matter?
		buf: make([]byte, 8*1024*1024),
		chunker: chunker.NewWithConfig(r, chunker.Pol(0x3DA3358B4DC173),
			chunker.ChunkerConfig{
				MinSize:     uint(math.Floor(float64(averageChunkSize) * minSizeMul)),
				MaxSize:     uint(math.Floor(float64(averageChunkSize) * maxSizeMul)),
				AverageBits: averageBits,
			}),
	}, nil
}

func (c *Roller) Roll() (fixity.Chunk, error) {
	// TODO(leeola): Add a peek method to break out of the loop if the end of the
	// roller is near. This way we don't create small tailing chunks if possible.

	chunk, err := c.chunker.Next(c.buf)
	// eof is okay to pass on immediately
	if err != nil {
		return fixity.Chunk{}, err
	}

	return fixity.Chunk{
		ChunkBytes: chunk.Data,
		Size:       int64(chunk.Length),
	}, nil
}
