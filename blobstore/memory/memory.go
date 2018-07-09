package memory

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"sync"

	base58 "github.com/jbenet/go-base58"
	"github.com/leeola/fixity"
	blake2b "github.com/minio/blake2b-simd"
)

// Store is a memory store used for testing.
type Store struct {
	mu sync.Mutex
	m  map[fixity.Ref][]byte
}

func New() *Store {
	return &Store{
		m: map[fixity.Ref][]byte{},
	}
}

func (s *Store) Read(_ context.Context, ref fixity.Ref) (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.m[ref]
	if !ok {
		return nil, os.ErrNotExist
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func (s *Store) Write(_ context.Context, b []byte) (fixity.Ref, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hB := blake2b.Sum256(b)
	ref := fixity.Ref(base58.Encode(hB[:]))
	s.m[ref] = b
	return ref, nil
}
