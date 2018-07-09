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

	"github.com/leeola/fixity"
)

const bsDir = "blobs"

// Blobstore implements a Fixity Blobstore for an simple Filesystem.
//
// NOTE: Blobstore is not safe for concurrent use out of process, but
// side effects are mostly harmless. Safe readers of partial writes
// should verify data regardless.
type Blobstore struct {
	mu   sync.Mutex
	path string
}

func New(path string) (*Blobstore, error) {
	if path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	path = filepath.Join(path, bsDir)

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return &Blobstore{
		path: path,
	}, nil
}

func (s *Blobstore) Read(_ context.Context, h fixity.Ref) (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if h == "" {
		return nil, errors.New("hash cannot be empty")
	}

	p := s.pathHash(string(h))

	rc, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	return rc, nil
}

func (s *Blobstore) Write(_ context.Context, b []byte) (fixity.Ref, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	h, err := fixity.Hash(b)
	if err != nil {
		return "", fmt.Errorf("hash: %v", err)
	}

	p := s.pathHash(string(h))

	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return "", fmt.Errorf("mkdirall: %v", err)
	}

	if err := ioutil.WriteFile(p, b, 0644); err != nil {
		return "", fmt.Errorf("writefile: %v", err)
	}

	return h, nil
}
