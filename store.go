package fixity

import (
	"context"
	"io"
)

type Store interface {
	Blob(ctx context.Context, ref Ref) (io.ReadCloser, error)
	Read(ctx context.Context, id string) (Mutation, Values, Reader, error)
	ReadRef(context.Context, Ref) (Mutation, Values, Reader, error)
	Write(ctx context.Context, id string, v Values, r io.Reader) ([]Ref, error)
	Querier
}
