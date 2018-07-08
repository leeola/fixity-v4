package nosign

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
	"github.com/leeola/fixity/chunk/resticfork"
	"github.com/leeola/fixity/reader/datareader"
	"github.com/leeola/fixity/util/wutil"
)

type Store struct {
	bs blobstore.ReadWriter
}

func New(bs blobstore.ReadWriter) (*Store, error) {
	return &Store{bs: bs}, nil
}

func (s *Store) Write(ctx context.Context, id string, v fixity.Values, r io.Reader) ([]fixity.Ref, error) {
	return s.WriteTime(ctx, time.Now(), id, v, r)
}

func (s *Store) WriteTime(ctx context.Context, t time.Time, id string, v fixity.Values, r io.Reader) ([]fixity.Ref, error) {
	// default to user namespace, ie ""
	return s.WriteTimeNamespace(ctx, t, id, "", v, r)
}

func (s *Store) WriteTimeNamespace(ctx context.Context,
	t time.Time, id, namespace string, v fixity.Values, r io.Reader) ([]fixity.Ref, error) {

	if v == nil && r == nil {
		return nil, errors.New("values and data cannot be nil")
	}

	var refs []fixity.Ref

	var dataRef fixity.Ref
	if r != nil {
		chunker, err := resticfork.New(r, resticfork.DefaultAverageChunkSize)
		if err != nil {
			return nil, fmt.Errorf("restic new: %v", err)
		}

		cHashes, totalSize, checksum, err := wutil.WriteChunks(ctx, s.bs, chunker)
		if err != nil {
			return nil, fmt.Errorf("writechunker: %v", err)
		}

		cHashes, err = wutil.WriteData(ctx, s.bs, cHashes, totalSize, checksum)
		if err != nil {
			return nil, fmt.Errorf("writecontent: %v", err)
		}

		dataRef = cHashes[len(cHashes)-1]
		refs = cHashes
	}

	var valuesRef fixity.Ref
	if v != nil {
		ref, err := wutil.WriteValues(ctx, s.bs, v)
		if err != nil {
			return nil, fmt.Errorf("writecontent: %v", err)
		}
		valuesRef = ref
		refs = append(refs, ref)
	}

	mutation := fixity.Mutation{
		Schema: fixity.Schema{
			SchemaType: fixity.BlobTypeMutation,
		},
		ID:        id,
		Time:      t,
		Data:      dataRef,
		ValuesMap: valuesRef,
	}

	ref, err := wutil.MarshalAndWrite(ctx, s.bs, mutation)
	if err != nil {
		return nil, fmt.Errorf("marshalandwrite mutation: %v", err)
	}

	return append(refs, ref), nil
}

func (s *Store) Blob(ctx context.Context, ref fixity.Ref) (io.ReadCloser, error) {
	rc, err := s.bs.Read(ctx, ref)
	if err != nil {
		// not wrapping to let error values fall through. The error context
		// from this store is likely meaningless here.
		return nil, err
	}

	return rc, nil
}

func (s *Store) Read(ctx context.Context, ref fixity.Ref) (fixity.Mutation, fixity.Values, io.Reader, error) {
	var mutation fixity.Mutation
	if err := blobstore.ReadAndUnmarshal(ctx, s.bs, ref, &mutation); err != nil {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("read mutation: %v", err)
	}

	if mutation.SchemaType != fixity.BlobTypeMutation {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("must read mutation blobs")
	}

	var values fixity.ValuesMap
	if mutation.ValuesMap != "" {
		if err := blobstore.ReadAndUnmarshal(ctx, s.bs, mutation.ValuesMap, &values); err != nil {
			return fixity.Mutation{}, nil, nil, fmt.Errorf("read values: %v", err)
		}
	}

	var data io.Reader
	if mutation.Data != "" {
		dr, err := datareader.New(ctx, s.bs, mutation.Data)
		if err != nil {
			return fixity.Mutation{}, nil, nil, fmt.Errorf("datareader new: %v", err)
		}
		data = dr
	}

	// values will be nil if not defined, which is okay.
	return mutation, values.Values, data, nil
}
