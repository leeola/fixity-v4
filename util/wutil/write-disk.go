package wutil

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/dchest/blake2b"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
	"github.com/leeola/fixity/chunk"
)

func WriteChunker(ctx context.Context, w blobstore.Writer, r chunk.Chunker) (
	refs []fixity.Ref, totalSize int64, contentHash string, err error) {

	hasher := blake2b.New256()

	var hashes []fixity.Ref
	for {
		c, err := r.Chunk(ctx)
		if err != nil && err != io.EOF {
			return nil, 0, "", fmt.Errorf("chunk: %v", err)
		}

		totalSize += c.Size

		if err == io.EOF {
			break
		}

		if _, err := hasher.Write(c.Bytes); err != nil {
			return nil, 0, "", fmt.Errorf("hasher write: %v", err)
		}

		h, err := w.Write(ctx, c.Bytes)
		if err != nil {
			return nil, 0, "", fmt.Errorf("blob write: %v", err)
		}

		hashes = append(hashes, h)
	}

	hash := hex.EncodeToString(hasher.Sum(nil)[:])
	return hashes, totalSize, hash, nil
}
