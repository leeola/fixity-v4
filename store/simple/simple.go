package simple

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
	blake2b "github.com/minio/blake2b-simd"
)

type Config struct {
	StorePath string
	Log       log15.Logger
}

type Simple struct {
	path string
	log  log15.Logger
}

func New(c Config) (*Simple, error) {
	if c.StorePath == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	return &Simple{
		log:  c.Log,
		path: c.StorePath,
	}, nil
}

func (s *Simple) Exists(h string) (bool, error) {
	p := filepath.Join(s.path, h)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "simple store failed to stat hash: %s", h)
	}
	return true, nil
}

func (s *Simple) Read(h string) (io.ReadCloser, error) {
	p := filepath.Join(s.path, h)
	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, store.HashNotFoundErr
	}
	return f, errors.Wrapf(err, "simple store failed to read hash: %s", h)
}

func (s *Simple) Hash(b []byte) string {
	h := blake2b.Sum256(b)
	return hex.EncodeToString(h[:])
}

func (s *Simple) Write(b []byte) (string, error) {
	h := s.Hash(b)
	if err := s.writeHash(h, b); err != nil {
		return "", err
	}
	return h, nil
}

func (s *Simple) WriteHash(h string, b []byte) error {
	expectedH := s.Hash(b)
	if h != expectedH {
		return store.HashNotMatchContentErr
	}
	return s.writeHash(h, b)
}

// writeHash is a trusted implementation of writeHash that does *not* verify the hash
//
// Verification of the content *must be done* before using this method to write.
func (s *Simple) writeHash(h string, b []byte) error {
	p := filepath.Join(s.path, h)
	err := ioutil.WriteFile(p, b, 0644)
	return errors.Wrap(err, "failed to write to disk")
}

func (s *Simple) List(max, offset int) (<-chan string, error) {
	return nil, errors.New("not implemented")
}

func (c Config) IsZero() bool {
	switch {
	case c.Log != nil:
		return false
	case c.StorePath != "":
		return false
	default:
		return true
	}
}
