package disk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
	base58 "github.com/jbenet/go-base58"
	blake2b "github.com/minio/blake2b-simd"
)

type Config struct {
	// Path is the *directory* to contain the store content and metadata.
	//
	// This will be created if it does not exist.
	Path string
}

// Disk implements a Fixity Store for an simple Filesystem.
type Disk struct {
	path string
	log  log15.Logger
}

func New(c Config) (*Disk, error) {
	if c.Path == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if err := os.MkdirAll(c.Path, 0755); err != nil {
		return nil, err
	}

	return &Disk{
		path: c.Path,
	}, nil
}

func (s *Disk) Read(_ context.Context, h string) (io.ReadCloser, error) {
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
