package disk

import (
	"bytes"
	"io"
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
	p := s.pathHash(h)

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

	p := s.pathHash(h)

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

func (s *Disk) Write(b []byte) (string, bool, error) {
	h := s.Hash(b)
	created, err := s.writeHash(h, b)
	if err != nil {
		return "", created, err
	}
	return h, created, nil
}

func (s *Disk) WriteHash(h string, b []byte) (bool, error) {
	expectedH := s.Hash(b)
	if h != expectedH {
		return false, fixity.ErrHashNotMatchBytes
	}
	return s.writeHash(h, b)
}

// writeHash is a trusted implementation of writeHash that does *not* verify the hash
//
// Verification of the content *must be done* before using this method to write.
func (s *Disk) writeHash(h string, b []byte) (bool, error) {
	p := s.pathHash(h)

	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return false, err
	}

	// Create or Truncate the file.
	//
	// By separating the truncate call, it allows us to know if the file was
	// created or not. We truncate to ensure we always write the requested data,
	// to prevent possible partial wrties and "bad data".
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil && !os.IsExist(err) {
		return false, err
	}
	// deferring the close is done below the truncate call.

	created := !os.IsExist(err)
	if !created {
		file, err := os.OpenFile(p, os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return false, err
		}
		f = file
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(b)); err != nil {
		return false, err
	}

	// Call Sync to ensure data wrote successfully.
	// Reference: https://joeshaw.org/dont-defer-close-on-writable-files/
	if err := f.Sync(); err != nil {
		return false, err
	}

	return created, nil
}

func (s *Disk) List() (<-chan string, error) {
	// TODO(leeola): Use a concurrent walking library to make this faster,
	// since Stdlib uses lexical order and we don't need deterministic results.

	ch := make(chan string)
	go func() {
		s.log.Debug("starting list walk")
		err := filepath.Walk(s.path, func(p string, _ os.FileInfo, _ error) error {
			// Trim the store path from the returned paths
			// TODO(leeola): remove the dirs from the h value
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
