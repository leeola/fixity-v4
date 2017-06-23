package disk

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
	base58 "github.com/jbenet/go-base58"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	blake2b "github.com/minio/blake2b-simd"
)

type Config struct {
	// Path is the *directory* to contain the store content and metadata.
	//
	// This will be created if it does not exist.
	Path string
	Log  log15.Logger
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

	if c.Log == nil {
		c.Log = log15.New()
	}

	return &Disk{
		log:  c.Log,
		path: c.Path,
	}, nil
}

func (s *Disk) Exists(h string) (bool, error) {
	p := filepath.Join(s.path, h)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "disk store failed to stat hash: %s", h)
	}
	return true, nil
}

func (s *Disk) Read(h string) (io.ReadCloser, error) {
	if h == "" {
		return nil, errors.New("hash cannot be empty")
	}

	p := filepath.Join(s.path, h)

	var rc io.ReadCloser
	rc, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, fixity.ErrHashNotFound
	}

	return rc, err
}

func (s *Disk) Hash(b []byte) string {
	hB := blake2b.Sum256(b)
	return base58.Encode(hB[:])
}

func (s *Disk) Write(b []byte) (string, error) {
	h := s.Hash(b)
	if err := s.writeHash(h, b); err != nil {
		return "", err
	}
	return h, nil
}

func (s *Disk) WriteHash(h string, b []byte) error {
	expectedH := s.Hash(b)
	if h != expectedH {
		return fixity.ErrHashNotMatchBytes
	}
	return s.writeHash(h, b)
}

// writeHash is a trusted implementation of writeHash that does *not* verify the hash
//
// Verification of the content *must be done* before using this method to write.
func (s *Disk) writeHash(h string, b []byte) error {
	p := filepath.Join(s.path, h)

	err := ioutil.WriteFile(p, b, 0644)
	return errors.Wrap(err, "failed to write to disk")
}

func (s *Disk) List() (<-chan string, error) {
	// TODO(leeola): Use a concurrent walking library to make this faster,
	// since Stdlib uses lexical order and we don't need deterministic results.

	ch := make(chan string)
	go func() {
		s.log.Debug("starting list walk")
		err := filepath.Walk(s.path, func(p string, _ os.FileInfo, _ error) error {
			// Trim the store path from the returned paths
			h, err := filepath.Rel(s.path, p)
			if err != nil {
				return err
			}

			// walk returns the base path, so ignore that
			if h == "." {
				return nil
			}

			ch <- h

			return nil
		})
		if err != nil {
			s.log.Error("list walk returned error", "err", err)
		}

		close(ch)
		s.log.Debug("done listing")
	}()

	return ch, nil
}
