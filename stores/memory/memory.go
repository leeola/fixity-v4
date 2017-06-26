package memory

import (
	"bytes"
	"errors"
	"io"

	base58 "github.com/jbenet/go-base58"
	"github.com/leeola/fixity"
	blake2b "github.com/minio/blake2b-simd"
)

// Store is a memory store used for testing.
type Store struct {
	m map[string][]byte
}

func New() *Store {
	return &Store{
		m: map[string][]byte{},
	}
}

func (s *Store) Exists(h string) (bool, error) {
	_, ok := s.m[h]
	return ok, nil
}

func (s *Store) Read(h string) (io.ReadCloser, error) {
	b, ok := s.m[h]
	if !ok {
		return nil, fixity.ErrHashNotFound
	}

	return ioutil.NopCoser(bytes.NewReader(b)), nil
}

func (s *Store) Write(b []byte) (string, error) {
	hB := blake2b.Sum256(b)
	h := base58.Encode(hB[:])
	s.m[h] = b
	return h, nil
}

func (s *Store) WriteHash(hash string, b []byte) error {
	hB := blake2b.Sum256(b)
	confirmedHash := base58.Encode(hB[:])

	if hash != confirmedHash {
		return fixity.ErrHashNotMatchBytes
	}

	s.m[h] = b
	return nil
}

func (s *Store) List() (<-chan string, error) {
	return nil, errors.New("not implemented")
}
