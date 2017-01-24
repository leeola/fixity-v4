package blev

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

type Config struct {
	BleveDir string

	// The store to rebuild the index from, if needed.
	Store store.Store

	// Optional.
	Log log15.Logger `json:"-"`
}

type Bleve struct {
	entryIndex   bleve.Index
	anchorIndex  bleve.Index
	indexDir     string
	indexVersion string
	log          log15.Logger
}

func New(c Config) (*Bleve, error) {
	if c.BleveDir == "" {
		return nil, errors.New("missing required Config field: BleveDir")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	if err := makeIndexDirs(c.BleveDir); err != nil {
		return nil, errors.Stack(err)
	}

	// open the entry index. If it does not exist, it will be created under Rebuild
	entryPath := filepath.Join(c.BleveDir, entryIndexDir)
	entryIndex, err := bleve.Open(entryPath)
	if err != nil && err != bleve.ErrorIndexMetaMissing {
		return nil, errors.Wrap(err, "failed to open or create bleve normal index")
	}

	// open the anchor index. If it does not exist, it will be created under Rebuild
	anchorPath := filepath.Join(c.BleveDir, anchorIndexDir)
	anchorIndex, err := bleve.Open(anchorPath)
	if err != nil && err != bleve.ErrorIndexMetaMissing {
		return nil, errors.Wrap(err, "failed to open or create bleve unique index")
	}

	indexVersion := LoadVersion(c.BleveDir)
	if indexVersion == "" {
		iv, err := NewVersion(c.BleveDir)
		if err != nil {
			return nil, errors.Stack(err)
		}
		indexVersion = iv
	}

	b := &Bleve{
		entryIndex:   entryIndex,
		anchorIndex:  anchorIndex,
		indexDir:     c.BleveDir,
		indexVersion: indexVersion,
		log:          c.Log,
	}

	return b, nil
}

func (b *Bleve) IndexVersion() string {
	return b.indexVersion
}

func makeIndexDirs(rootPath string) error {
	entryPath := filepath.Join(rootPath, entryIndexDir)
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		if err := os.Mkdir(entryPath, 0755); err != nil {
			return errors.Stack(err)
		}
	}

	anchorPath := filepath.Join(rootPath, anchorIndexDir)
	if _, err := os.Stat(anchorPath); os.IsNotExist(err) {
		if err := os.Mkdir(anchorPath, 0755); err != nil {
			return errors.Stack(err)
		}
	}

	return nil
}

func removeIndexDirs(rootPath string) error {
	entryPath := filepath.Join(rootPath, entryIndexDir)
	if err := os.RemoveAll(entryPath); err != nil {
		return errors.Stack(err)
	}

	anchorPath := filepath.Join(rootPath, anchorIndexDir)
	if err := os.RemoveAll(anchorPath); err != nil {
		return errors.Stack(err)
	}

	return nil
}

func NewVersion(indexDir string) (string, error) {
	indexVersion := strconv.Itoa(int(time.Now().Unix()))
	if err := SaveVersion(indexDir, indexVersion); err != nil {
		return "", err
	}
	return indexVersion, nil
}

func SaveVersion(indexDir, version string) error {
	err := ioutil.WriteFile(filepath.Join(indexDir, "version"), []byte(version), 0644)
	return errors.Wrap(err, "failed to save version")
}

func LoadVersion(indexDir string) string {
	v, _ := ioutil.ReadFile(filepath.Join(indexDir, "version"))
	return string(v)
}
