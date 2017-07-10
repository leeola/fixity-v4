package sync

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/leeola/fixity"
)

type Config struct {
	Path      string
	Folder    string
	Recursive bool
	Fixity    fixity.Fixity
}

type Sync struct {
	config Config
	fixi   fixity.Fixity
	c      chan string
}

func New(c Config) (*Sync, error) {
	if c.Path == "" {
		return nil, errors.New("missing reqired config: Path")
	}

	if c.Fixity == nil {
		return nil, errors.New("missing reqired config: Fixity")
	}

	if c.Folder == "" {
		c.Folder = filepath.Dir(c.Path)
	}

	if c.Folder == "" {
		return nil, errors.New("at least one folder is required")
	}

	return &Sync{
		config: c,
		fixi:   c.Fixity,
	}, nil
}

func (s *Sync) Sync() error {
	defer func() {
		if s.c != nil {
			close(s.c)
			s.c = nil
		}
	}()

	fi, err := os.Stat(s.config.Path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return s.syncDir(s.config.Path)
	} else {
		return s.syncFile(s.config.Path)
	}
}

func (s *Sync) syncDir(path string) error {
	return errors.New("not implemented")
}

func (s *Sync) replaceFile(path string, c fixity.Content) error {
	// using O_CREATE just to be safe, in case something external deletes the
	// file, no reason we can't still create it.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	rc, err := c.Read()
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return err
	}

	return f.Sync()
}

func (s *Sync) syncFile(path string) error {
	// update the update chan with the latest file sync'd
	defer s.update(path)

	c, err := s.uploadFile(s.config.Path)
	if err != nil {
		return err
	}

	switch c.Index {
	case 1:
		// if the index is 1, this content was appended and was not duplicate.
		// Syncing back to the filesystem is not needed, so append it.
		return nil
	case 0:
		// if the index is 0, we cannot assert if the file needs to be sync'd
		// or not. Return an error.
		//
		// This ensures in the event that we don't know the file order,
		// we don't overwrite users files.
		return errors.New("syncFile: unable to sync, unknown Content index of 0")
	}

	// if the index was larger than 1, then it's either unknown or an older blob.
	// In that case, read the file from fixity and write to disk.
	return s.replaceFile(path, c)
}

func (s *Sync) update(path string) {
	if s.c != nil {
		s.c <- path
	}
}

func (s *Sync) Updates() <-chan string {
	if s.c == nil {
		s.c = make(chan string, 50)
	}
	return s.c
}

func (s *Sync) uploadFile(path string) (fixity.Content, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return fixity.Content{}, err
	}
	defer f.Close()

	// TODO(leeola): include unix metadata
	req := fixity.NewWrite(path, f)
	req.IgnoreDuplicateBlob = true

	return s.fixi.WriteRequest(req)
}
