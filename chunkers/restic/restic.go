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

// min is pretty harmless so we should aim for it to be close to the original
// value.
const minSizeMul = 0.9

// restic chunker puts a hard cap on the max size so this should probably be
// quite a bit larger than the desired max
const maxSizeMul = 2.0

type Chunker struct {
	buf     []byte
	chunker *chunker.Chunker
}

func New(r io.Reader, averageChunkSize uint64) (*Chunker, error) {
	if r == nil {
		return nil, errors.New("missing Reader")
	}

	min := uint(math.Floor(float64(averageChunkSize) * minSizeMul))
	max := uint(math.Floor(float64(averageChunkSize) * maxSizeMul))

	return &Chunker{
		// does this size matter?
		buf: make([]byte, 8*1024*1024),
		chunker: chunker.NewWithConfig(r, chunker.Pol(0x3DA3358B4DC173),
			chunker.ChunkerConfig{
				MinSize:     min,
				MaxSize:     max,
				AverageBits: averageBits,
			}),
	}, nil
}

func (c *Chunker) Chunk() (fixity.Chunk, error) {

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
