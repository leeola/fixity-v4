package sync

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/leeola/fixity"
)

type Iter interface {
	Next() (iterHasValue bool)
	Value() (c fixity.Content, err error)
}

type Config struct {
	Path      string
	Folder    string
	Recursive bool
	Fixity    fixity.Fixity
}

type Sync struct {
	config Config
	fixi   fixity.Fixity

	ch  chan walkResult
	c   fixity.Content
	err error
}

type walkResult struct {
	Path string
	Err  error
}

func New(c Config) (*Sync, error) {
	if c.Path == "" {
		return nil, errors.New("missing reqired config: Path")
	}

	if c.Fixity == nil {
		return nil, errors.New("missing reqired config: Fixity")
	}

	// if folder is still empty, check if the Path is a directory.
	// this supports the `sync ./dir` usage.
	if c.Folder == "" {
		fi, err := os.Stat(c.Path)
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			c.Folder = filepath.Base(c.Path)
		} else if p := filepath.Base(filepath.Dir(c.Path)); p != "." {
			c.Folder = p
		}
	}

	// enforcing relative folders will allow exporting to be a bit easier/safer.
	if filepath.IsAbs(c.Folder) {
		return nil, errors.New("folder must be relative")
	}

	if c.Folder == "" {
		return nil, errors.New("at least one folder is required")
	}

	return &Sync{
		config: c,
		fixi:   c.Fixity,
	}, nil
}

func (s *Sync) walk() {
	err := filepath.Walk(s.config.Path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if s.config.Recursive || path == s.config.Path {
				return nil
			} else {
				return filepath.SkipDir
			}
		}

		s.ch <- walkResult{Path: path}

		return nil
	})
	if err != nil {
		s.ch <- walkResult{Err: err}
	}
	close(s.ch)
}

func (s *Sync) Next() bool {
	if s.ch == nil {
		s.ch = make(chan walkResult)
		go s.walk()
	}

	walkResult, ok := <-s.ch
	if !ok {
		return false
	}

	if walkResult.Err != nil {
		s.c = fixity.Content{}
		s.err = walkResult.Err
		// return true because there is an error that the caller
		// should grab via .Value()
		return true
	}

	s.c, s.err = s.syncFile(walkResult.Path)
	return true
}

func (s *Sync) Value() (fixity.Content, error) {
	return s.c, s.err
}

func (s *Sync) replaceFile(path string, outdated fixity.Content) error {
	// using O_CREATE just to be safe, in case something external deletes the
	// file, no reason we can't still create it.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c, err := s.fixi.Read(outdated.Id)
	if err != nil {
		return err
	}

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

func (s *Sync) syncFile(path string) (fixity.Content, error) {
	c, err := s.uploadFile(path)
	if err != nil {
		return fixity.Content{}, err
	}

	switch c.Index {
	case 1:
		// if the index is 1, this content was appended and was not duplicate.
		// Syncing back to the filesystem is not needed, so append it.
		return c, nil
	case 0:
		// if the index is 0, we cannot assert if the file needs to be sync'd
		// or not. Return an error.
		//
		// This ensures in the event that we don't know the file order,
		// we don't overwrite users files.
		return fixity.Content{}, errors.New("syncFile: unable to sync, unknown Content index of 0")
	}

	// if the index was larger than 1, then it's either unknown or an older blob.
	// In that case, read the file from fixity and write to disk.
	if err := s.replaceFile(path, c); err != nil {
		return fixity.Content{}, err
	}

	return c, nil
}

func (s *Sync) uploadFile(path string) (fixity.Content, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return fixity.Content{}, err
	}
	defer f.Close()

	// by resolving the path relative to the folder, and then joining
	// them, we ensure the id is always a subdirectory file of the c.Folder.
	// While also ensuring we don't double up on the root folder.
	// Eg:
	//    sync foodir
	// doesn't become
	//    sync foodir/foodir/foofile
	id, err := filepath.Rel(s.config.Folder, path)
	if err != nil {
		return fixity.Content{}, err
	}
	id = filepath.Join(s.config.Folder, id)

	// TODO(leeola): include unix metadata
	req := fixity.NewWrite(id, f)
	req.IgnoreDuplicateBlob = true

	return s.fixi.WriteRequest(req)
}
