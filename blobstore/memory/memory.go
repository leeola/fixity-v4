package memory

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"sync"

	base58 "github.com/jbenet/go-base58"
	blake2b "github.com/minio/blake2b-simd"
)

// Store is a memory store used for testing.
type Store struct {
	mu sync.Mutex
	m  map[string][]byte
}

func New() *Store {
	return &Store{
		m: map[string][]byte{},
	}
}

func (s *Store) Read(_ context.Context, h string) (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.m[h]
	if !ok {
		return nil, os.ErrNotExist
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func (s *Store) Write(_ context.Context, b []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hB := blake2b.Sum256(b)
	h := base58.Encode(hB[:])
	s.m[h] = b
	return h, nil
}
