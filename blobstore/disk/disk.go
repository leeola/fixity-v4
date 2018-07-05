package disk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	base58 "github.com/jbenet/go-base58"
	blake2b "github.com/minio/blake2b-simd"
)

// Disk implements a Fixity Store for an simple Filesystem.
//
// NOTE: Disk is not safe for concurrent use out of process, but
// side effects are mostly harmless. Safe readers of partial writes
// should verify data regardless.
type Disk struct {
	mu   sync.Mutex
	path string
}

func New(path string) (*Disk, error) {
	if path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return &Disk{
		path: path,
	}, nil
}

func (s *Disk) Read(_ context.Context, h string) (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if h == "" {
		return nil, errors.New("hash cannot be empty")
	}

	p := s.pathHash(h)

	rc, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	return rc, nil
}

func (s *Disk) Hash(b []byte) string {
	hB := blake2b.Sum256(b)
	return base58.Encode(hB[:])
}

func (s *Disk) Write(_ context.Context, b []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	h := s.Hash(b)
	p := s.pathHash(h)

	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return "", fmt.Errorf("mkdirall: %v", err)
	}

	if err := ioutil.WriteFile(p, b, 0644); err != nil {
		return "", fmt.Errorf("writefile: %v", err)
	}

	return h, nil
}
