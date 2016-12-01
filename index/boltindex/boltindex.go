package boltindex

import (
	"encoding/binary"
	"path/filepath"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

const (
	dbFilename = "boltindex.db"
)

var (
	// the bucket name used to store map[int]string where int is the Nth hash
	// added to the store and string is the hash.
	entryBucketName = []byte("entry")

	// the bucket name used to store data about the data. Ie, what version the
	// index is, etc.
	metaBucketName      = []byte("meta")
	metaEntryCountKey   = []byte("entryCount")
	metaIndexVersionKey = []byte("indexVersion")
)

type Config struct {
	// The path to the directory containing the Boltdb database file.
	BoltPath string

	// Indexes currently get their updates by acting as middleware between the
	// real Store and the Node. Due to this, BoltIndex must implement Index *and*
	// Store, and the actual Store must be given to BoltIndex so it can be wrapped.
	Store store.Store `toml:"-"`
}

type BoltIndex struct {
	version    string
	entryCount int
	store      store.Store

	// the db where we actually store the index information
	db *bolt.DB
}

func New(c Config) (*BoltIndex, error) {
	if c.BoltPath == "" {
		return nil, errors.New("missing required Config field: BoltPath")
	}
	if c.Store == nil {
		return nil, errors.New("missing required Config field: Store")
	}

	db, err := bolt.Open(filepath.Join(c.BoltPath, dbFilename), 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bolt index")
	}

	bucketNames := [][]byte{entryBucketName, metaBucketName}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bName := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(bName); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create buckets")
	}

	bi := &BoltIndex{
		store: c.Store,
		db:    db,
	}

	if err := bi.initMeta(); err != nil {
		return nil, err
	}

	return bi, nil
}

func (bi *BoltIndex) initMeta() error {
	return bi.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(metaBucketName)
		v := b.Get(metaIndexVersionKey)
		if v != nil {
			bi.version = string(v)
		} else {
			// TODO(leeola): use hostname or uuid for the index versioning.
			bi.version = strconv.FormatInt(time.Now().Unix(), 10)
			if err := b.Put(metaIndexVersionKey, []byte(bi.version)); err != nil {
				return err
			}
		}

		if v := b.Get(metaEntryCountKey); v != nil {
			bi.entryCount = btoi(v)
		}
		return nil
	})
}

func (bi *BoltIndex) GetInt(bucket, key []byte) (int, error) {
	var i int
	err := bi.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(key)
		if v != nil {
			i = btoi(v)
		}
		return nil
	})
	return i, err
}

func (bi *BoltIndex) GetString(bucket, key []byte) (string, error) {
	var s string
	err := bi.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(key)
		if v != nil {
			s = string(v)
		}
		return nil
	})
	return s, err
}

func (bi *BoltIndex) GetEntry(i int) (string, error) {
	return bi.GetString(entryBucketName, itob(i))
}

func itob(i int) []byte {
	b := make([]byte, 32)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
