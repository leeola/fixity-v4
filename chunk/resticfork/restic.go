// Restic Fork uses a fork of the restic chunking codebase which supports
// custom chunk sizing.
//
// This is likely being deprecated soon, in favor of standard restic chunking.

package resticfork

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/leeola/chunker"
	"github.com/leeola/errors"
	"github.com/leeola/fixity/chunk"
)

const (
	// DefaultAverageChunkSize that the chunker will attempt to produce,
	// if given.
	//
	// See WriteRequest.AverageChunkSize documentation for further explanation
	// of chunk sizes.
	//
	// This value is 1024^2*1 bytes, 1MiB, chosen as a compromise for the
	// small and large files.
	//
	// Ideally it should strike a balance between too many chunks written while
	// still mitigating data duplication caused by frequent changes to files
	// and written data.
	//
	// If written data is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	DefaultAverageChunkSize = medChunkSize

	// medChunkSize is the size SetChunkSizeFromX will use for medium data.
	//
	// This value is 1024^2 bytes, 1MiB, chosen as a balance for medium data.
	//
	// Note that using this medChunkSize is decided via the medCutoff, below.
	medChunkSize uint64 = 1048576

	// averageBits of 14 is resulting in roughly chunks of 12KiB by default.
	//
	// This size should be low enough to allow small writes to be split up,
	// even at the cost of write performance (within reason).
	averageBits = 14

	// min is pretty harmless so we should aim for it to be close to the original
	// value.
	minSizeMul = 0.9

	// restic chunker puts a hard cap on the max size so this should probably be
	// quite a bit larger than the desired max.
	//
	// This is currently set to 3x. So if restic is unable to find a chunk
	// boundry within 3x of the average chunk size, restic will force a
	// chunkboundry.
	// This prevents a boundry from never being found and a potential chunk
	// size of GiBs/etc having to be stored, sent, etc.
	//
	// Fixity expects reasonably small chunk sizes, but we forcing chunk sizes
	// is also very bad as it's not a deterministic boundry, as it would
	// depend on the last boundry. So, this value should strike a good
	// balance between those two issues.
	maxSizeMul = 3.0
)

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

func (c *Chunker) Chunk(_ context.Context) (chunk.Chunk, error) {

	// TODO(leeola): Add a peek method to break out of the loop if the end of the
	// roller is near. This way we don't create small tailing chunks if possible.

	ch, err := c.chunker.Next(c.buf)
	if err == io.EOF {
		return chunk.Chunk{}, err
	}
	if err != nil {
		return chunk.Chunk{}, fmt.Errorf("next: %v", err)
	}

	return chunk.Chunk{
		Bytes: ch.Data,
		Size:  int64(ch.Length),
	}, nil
}
