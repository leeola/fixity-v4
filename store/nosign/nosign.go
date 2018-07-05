package nosign

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
	"github.com/leeola/fixity/chunk/resticfork"
	"github.com/leeola/fixity/util/wutil"
)

type Store struct {
	bs blobstore.ReadWriter
}

func New(bs blobstore.ReadWriter) (*Store, error) {
	return &Store{bs: bs}, nil
}

func (s *Store) Write(ctx context.Context, id string, r io.Reader) ([]fixity.Ref, error) {
	return s.WriteTime(ctx, time.Now(), id, r)
}

func (s *Store) WriteTime(ctx context.Context, t time.Time, id string, r io.Reader) ([]fixity.Ref, error) {
	chunker, err := resticfork.New(r, resticfork.DefaultAverageChunkSize)
	if err != nil {
		return nil, fmt.Errorf("restic new: %v", err)
	}

	cHashes, _, _, err := wutil.WriteChunker(ctx, s.bs, chunker)
	if err != nil {
		return nil, fmt.Errorf("writechunker: %v", err)
	}

	return cHashes, nil
}
