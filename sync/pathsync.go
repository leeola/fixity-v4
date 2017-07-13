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

// TODO(leeola): provide a store path required field, to help ensure Fixity
// can never upload it's own store and loop endlessly.
type Config struct {
	Path      string
	Folder    string
	Recursive bool
	Fixity    fixity.Fixity
}

type Sync struct {
	config Config
	fixi   fixity.Fixity

	trimPath, path, folder string

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

	trimPath, path, folder, err := ResolveDirs(c.Path, c.Folder)
	if err != nil {
		return nil, err
	}

	if folder == "" {
		return nil, errors.New("at least one folder is required")
	}

	return &Sync{
		config:   c,
		fixi:     c.Fixity,
		trimPath: trimPath,
		path:     path,
		folder:   folder,
	}, nil
}

func (s *Sync) walk() {
	err := filepath.Walk(s.path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if s.config.Recursive || path == s.path {
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

	// by resolving the path relative to the trimPath, and then joining
	// them, we ensure the id is always a subdirectory file of the c.Folder.
	// While also ensuring we don't double up on the root folder.
	// Eg:
	//    sync foodir
	// doesn't become
	//    sync foodir/foodir/foofile
	// which is
	//    sync <providedDir>/<filePath>
	//
	// Much of the logic for this is provided via ResolveDirs
	id, err := filepath.Rel(s.trimPath, path)
	if err != nil {
		return fixity.Content{}, err
	}
	id = filepath.Join(s.folder, id)

	// TODO(leeola): include unix metadata
	req := fixity.NewWrite(id, f)
	req.IgnoreDuplicateBlob = true

	return s.fixi.WriteRequest(req)
}

func ResolveDirs(p, explicitFolder string) (trimPath, path, folder string, err error) {
	p, err = filepath.Abs(p)
	if err != nil {
		return "", "", "", err
	}

	fi, err := os.Stat(p)
	if err != nil {
		return "", "", "", err
	}

	var dirPath, fileName string
	if fi.IsDir() {
		dirPath = p
	} else {
		dirPath = filepath.Dir(p)
		fileName = filepath.Base(p)
	}

	return resolveDirs(dirPath, fileName, explicitFolder)
}

func resolveDirs(dirPath, fileName, explicitFolder string) (trimPath, path, folder string, err error) {
	if dirPath == "" {
		return "", "", "", errors.New("resolveDirs: directory is required")
	}
	if !filepath.IsAbs(dirPath) {
		return "", "", "", errors.New("resolveDirs: must provide absolute dir")
	}
	if filepath.IsAbs(explicitFolder) {
		return "", "", "", errors.New("resolveDirs: folder cannot be absolute")
	}

	if explicitFolder != "" {
		folder = explicitFolder
	} else {
		base := filepath.Base(dirPath)
		if base == "/" {
			return "", "", "", errors.New(
				"resolveDirs: must provid folder if no available directory to assert folder from")
		}
		// this should never happen, but worth checking
		if base == "." {
			return "", "", "", errors.New("resolveDirs: base resolved to '.'")
		}
		folder = base
	}

	// TODO(leeola): figure out what to do if a sole dir is provided *and* the
	// folder is provided. Eg, do we want to embed the dir in the folder? Or
	// ignore it, and put the dir's files in the providedFolder? etc.

	return dirPath, filepath.Join(dirPath, fileName), folder, nil
}
