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

	cHashes, totalSize, checksum, err := wutil.WriteChunks(ctx, s.bs, chunker)
	if err != nil {
		return nil, fmt.Errorf("writechunker: %v", err)
	}

	cHashes, err = wutil.WriteContent(ctx, s.bs, cHashes, totalSize, checksum)
	if err != nil {
		return nil, fmt.Errorf("writecontent: %v", err)
	}

	mutation := fixity.Mutation{
		ID:      id,
		Time:    t.String(), // TODO(leeola): parse?
		Content: cHashes[len(cHashes)-1],
	}

	ref, err := wutil.MarshalAndWrite(ctx, s.bs, mutation)
	if err != nil {
		return nil, fmt.Errorf("marshalandwrite mutation: %v", err)
	}

	return append(cHashes, ref), nil
}
