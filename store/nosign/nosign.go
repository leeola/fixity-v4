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
	"github.com/leeola/fixity/config"
	"github.com/leeola/fixity/index"
	"github.com/leeola/fixity/q"
	"github.com/leeola/fixity/reader/datareader"
	"github.com/leeola/fixity/util/wutil"
	"github.com/leeola/fixity/value"
)

type Config struct {
	BlobstoreKey string
	IndexKey     string
}

type Store struct {
	// embedded because the store exposes the same methods.
	index.Querier

	bstor fixity.ReadWriter
	index index.Indexer
}

func New(name string, fc config.Config) (*Store, error) {
	var c Config
	if err := fc.StoreConfig(name, &c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	bs, err := fixity.NewBlobstoreFromConfig(c.BlobstoreKey, fc)
	if err != nil {
		return nil, fmt.Errorf("blobstoreFromConfig: %v", err)
	}

	ix, err := fixity.NewIndexFromConfig(c.IndexKey, fc)
	if err != nil {
		return nil, fmt.Errorf("blobstoreFromConfig: %v", err)
	}

	return &Store{bstor: bs, index: ix, Querier: ix}, nil
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

	var (
		data    *fixity.DataSchema
		dataRef fixity.Ref
	)
	if r != nil {
		chunker, err := resticfork.New(r, resticfork.DefaultAverageChunkSize)
		if err != nil {
			return nil, fmt.Errorf("restic new: %v", err)
		}

		cHashes, totalSize, checksum, err := wutil.WriteChunks(ctx, s.bstor, chunker)
		if err != nil {
			return nil, fmt.Errorf("writechunker: %v", err)
		}

		cHashes, d, err := wutil.WriteData(ctx, s.bstor, cHashes, totalSize, checksum)
		if err != nil {
			return nil, fmt.Errorf("writecontent: %v", err)
		}
		data = d
		dataRef = cHashes[len(cHashes)-1]
		refs = cHashes
	}

	var valuesRef fixity.Ref
	if v != nil {
		ref, err := wutil.WriteValues(ctx, s.bstor, v)
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
		ID:           id,
		Time:         t,
		DataSchema:   dataRef,
		ValuesSchema: valuesRef,
	}

	ref, err := wutil.MarshalAndWrite(ctx, s.bstor, mutation)
	if err != nil {
		return nil, fmt.Errorf("marshalandwrite mutation: %v", err)
	}

	if err := s.index.Index(ref, mutation, data, v); err != nil {
		return nil, fmt.Errorf("index: %v", err)
	}

	return append(refs, ref), nil
}

func (s *Store) Blob(ctx context.Context, ref fixity.Ref) (io.ReadCloser, error) {
	rc, err := s.bstor.Read(ctx, ref)
	if err != nil {
		// not wrapping to let error values fall through. The error context
		// from this store is likely meaningless here.
		return nil, err
	}

	return rc, nil
}

func (s *Store) Read(ctx context.Context, id string) (
	fixity.Mutation, fixity.Values, fixity.Reader, error) {

	matches, err := s.Query(q.New().Eq(index.FIDKey, value.String(id)))
	if err != nil {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("query id: %v", err)
	}

	matchesLen := len(matches)
	tooManyMatches := matchesLen > 1
	noMatches := matchesLen == 0

	if tooManyMatches {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("id matched more than once")
	}

	if noMatches {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("id not found")
	}

	return s.ReadRef(ctx, matches[0].Ref)
}

func (s *Store) ReadRef(ctx context.Context, ref fixity.Ref) (
	fixity.Mutation, fixity.Values, fixity.Reader, error) {

	var mutation fixity.Mutation
	if err := blobstore.ReadAndUnmarshal(ctx, s.bstor, ref, &mutation); err != nil {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("read mutation: %v", err)
	}

	if mutation.SchemaType != fixity.BlobTypeMutation {
		return fixity.Mutation{}, nil, nil, fmt.Errorf("must read mutation blobs")
	}

	var values fixity.ValuesSchema
	if mutation.ValuesSchema != "" {
		if err := blobstore.ReadAndUnmarshal(ctx, s.bstor, mutation.ValuesSchema, &values); err != nil {
			return fixity.Mutation{}, nil, nil, fmt.Errorf("read values: %v", err)
		}
	}

	var data fixity.Reader
	if mutation.DataSchema != "" {
		dr, err := datareader.New(ctx, s.bstor, mutation.DataSchema)
		if err != nil {
			return fixity.Mutation{}, nil, nil, fmt.Errorf("datareader new: %v", err)
		}
		data = dr
	}

	// values will be nil if not defined, which is okay.
	return mutation, values.Values, data, nil
}
