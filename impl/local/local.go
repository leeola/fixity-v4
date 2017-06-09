package local

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	fixi "github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
)

var (
	blockMetaBucketKey = []byte("blockMeta")
	lastBlockKey       = []byte("lastBlock")
)

type Config struct {
	Index    fixity.Index `toml:"-"`
	Store    fixity.Store `toml:"-"`
	Log      log15.Logger `toml:"-"`
	RootPath string       `toml:"rootPath"`
}

type Local struct {
	config Config
	db     *bolt.DB
	index  fixity.Index
	store  fixity.Store
	log    log15.Logger
}

func New(c Config) (*Local, error) {
	if c.RootPath == "" {
		return nil, errors.New("missing required config: rootPath")
	}

	if c.Index == nil {
		return nil, errors.New("missing reqired config: Index")
	}

	if c.Store == nil {
		return nil, errors.New("missing reqired config: Store")
	}

	if c.Log == nil {
		c.Log = log15.New()
	}

	dbPath := filepath.Join(c.RootPath, "local", "blocks.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Local{
		config: c,
		db:     db,
		index:  c.Index,
		store:  c.Store,
		log:    c.Log,
	}, nil
}

func (l *Local) Blob(h string) ([]byte, error) {
	rc, err := l.store.Read(h)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (l *Local) Search(q *q.Query) ([]string, error) {
	return l.index.Search(q)
}

func (l *Local) getHead() (string, error) {
	var h string
	err := l.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockMetaBucketKey)
		// if bucket does not exist, this will be nil
		if bkt == nil {
			return nil
		}

		hB := bkt.Get(lastBlockKey)
		if hB != nil {
			h = string(hB)
		}

		return nil
	})
	return h, err
}

func (l *Local) setHead(h string) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(blockMetaBucketKey)
		if err != nil {
			return err
		}

		return bkt.Put(lastBlockKey, []byte(h))
	})
}

func (l *Local) Create(id string, r io.Reader, f ...fixi.Field) ([]string, error) {
	if r == nil {
		return nil, errors.New("no data given to write")
	}

	var hashes []string

	// TODO(leeola): roll reader
	rHash, err := WriteReader(l.store, r)
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, rHash)

	previousBlockHash, err := l.getHead()
	if err != nil {
		return nil, err
	}

	if id != "" {
		l.log.Warn("loading previous content hash not implemented")
	}

	content := fixi.Content{
		Id:            id,
		BlobHash:      rHash,
		IndexedFields: f,
	}

	cHash, err := MarshalAndWrite(l.store, content)
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, cHash)

	var lastBlock int
	// Get the previous block hash and count
	if previousBlockHash != "" {
		var prevBlock fixi.Block
		if err := ReadAndUnmarshal(l.store, previousBlockHash, &prevBlock); err != nil {
			return nil, err
		}
		lastBlock = prevBlock.Block
	}

	block := fixi.Block{
		// zero value is okay for both of these.
		Block:             lastBlock + 1,
		PreviousBlockHash: previousBlockHash,
		ContentHash:       cHash,
	}

	bHash, err := MarshalAndWrite(l.store, block)
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, bHash)

	return hashes, nil
}

func (l *Local) Delete(h string) error {
	return errors.New("not implemented")
}

func (l *Local) Update(h string, r io.Reader, f ...fixi.Field) ([]string, error) {
	return nil, errors.New("not implemented")
}

// WriteReader writes the given reader's content to the store.
func WriteReader(s fixity.Store, r io.Reader) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if r == nil {
		return "", errors.New("Reader is nil")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to readall")
	}

	h, err := s.Write(b)
	return h, errors.Wrap(err, "store failed to write")
}

// MarshalAndWrite marshals the given interface to json and writes that to the store.
func MarshalAndWrite(s fixity.Store, v interface{}) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if v == nil {
		return "", errors.New("Interface is nil")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func ReadAll(s fixity.Store, h string) ([]byte, error) {
	rc, err := s.Read(h)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func ReadAndUnmarshal(s fixity.Store, h string, v interface{}) error {
	_, err := ReadAndUnmarshalWithBytes(s, h, v)
	return err
}

func ReadAndUnmarshalWithBytes(s fixity.Store, h string, v interface{}) ([]byte, error) {
	b, err := ReadAll(s, h)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return nil, errors.Stack(err)
	}

	return b, nil
}
