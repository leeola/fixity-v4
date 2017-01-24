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
	"github.com/leeola/kala/store/crypto"
	blake2b "github.com/minio/blake2b-simd"
)

type Config struct {
	StorePath string
	Log       log15.Logger

	// Embedded so the config.toml user can write the values in the same object as
	// the store config.
	crypto.Config
}

type Simple struct {
	path string
	log  log15.Logger

	// Optional! Will be nil.
	cryptoer crypto.Cryptoer
}

func New(c Config) (*Simple, error) {
	if c.StorePath == "" {
		return nil, errors.New("missing required Config field: Path")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	// If there is no crypto configured by the user, do nothing.
	var cryptoer crypto.Cryptoer
	if c.UsesCrypto() {
		cer, err := crypto.New(c.Config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create cryptoer")
		}
		cryptoer = cer
	}

	return &Simple{
		log:      c.Log,
		path:     c.StorePath,
		cryptoer: cryptoer,
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

	var rc io.ReadCloser
	rc, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, store.HashNotFoundErr
	}

	if s.cryptoer != nil {
		dRc, err := crypto.DecryptReadCloser(s.cryptoer, rc)
		if err != nil {
			return nil, errors.Stack(err)
		}
		rc = dRc
	}

	return rc, errors.Wrapf(err, "simple store failed to read hash: %s", h)
}

func (s *Simple) Hash(b []byte) string {
	h := blake2b.Sum256(b)
	return "blake2b-" + hex.EncodeToString(h[:])
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

	// TODO(leeola): only write to disk if data does not exist.

	if s.cryptoer != nil {
		eb, err := s.cryptoer.Encrypt(b)
		if err != nil {
			return errors.Stack(err)
		}
		b = eb
	}

	err := ioutil.WriteFile(p, b, 0644)
	return errors.Wrap(err, "failed to write to disk")
}

func (s *Simple) List() (<-chan string, error) {
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
