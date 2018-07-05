package nosign

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
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

func (s *Store) WriteTime(_ context.Context, t time.Time, id string, r io.Reader) ([]fixity.Ref, error) {
	return nil, errors.New("not implemented")
}
