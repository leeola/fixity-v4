package bolt

import (
	"encoding/binary"
	"path/filepath"
	"time"

	boltdb "github.com/boltdb/bolt"
	"github.com/leeola/errors"
	"github.com/leeola/kala/database"
)

const (
	dbFilename = "bolt.db"
)

// Buckets and keys
var (
	// The name of the bucket for this nodes metadata
	nodeBucketName = []byte("node")
	nodeIdKey      = []byte("nodeId")

	// The name of the bucket for peer related data
	peersBucketName = []byte("peers")

	// the bucket name used to store map[int]string where int is the Nth hash
	// added to the store and string is the hash.
	indexEntryBucketName = []byte("indexEntry")

	// the bucket name used to store data about the data. Ie, what version the
	// index is, etc.
	indexMetaBucketName      = []byte("indexMeta")
	indexMetaEntryCountKey   = []byte("entryCount")
	indexMetaIndexVersionKey = []byte("indexVersion")

	// the bucket used to store random metadata
	metadataBucketName = []byte("metadata")
)

type Config struct {
	// The path to the directory containing the Boltdb database file.
	BoltPath string
}

type Bolt struct {
	// the db where we actually store the index information
	db *boltdb.DB
}

func New(c Config) (*Bolt, error) {
	if c.BoltPath == "" {
		return nil, errors.New("missing required Config field: BoltPath")
	}

	db, err := boltdb.Open(filepath.Join(c.BoltPath, dbFilename), 0600, &boltdb.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bolt index")
	}

	bucketNames := [][]byte{
		nodeBucketName,
		indexEntryBucketName,
		indexMetaBucketName,
		peersBucketName,
		metadataBucketName,
	}
	err = db.Update(func(tx *boltdb.Tx) error {
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

	b := &Bolt{
		db: db,
	}

	return b, nil
}

func (b *Bolt) GetInt(bucket, key []byte) (int, error) {
	var i int
	err := b.db.View(func(tx *boltdb.Tx) error {
		bkt := tx.Bucket(bucket)
		v := bkt.Get(key)
		if v == nil {
			return database.ErrNoRecord
		}
		i = btoi(v)
		return nil
	})
	return i, err
}

func (b *Bolt) GetString(bucket, key []byte) (string, error) {
	var s string
	err := b.db.View(func(tx *boltdb.Tx) error {
		bkt := tx.Bucket(bucket)
		v := bkt.Get(key)
		if v == nil {
			return database.ErrNoRecord
		}
		s = string(v)
		return nil
	})
	return s, err
}

func (b *Bolt) GetEntry(i int) (string, error) {
	return b.GetString(indexEntryBucketName, itob(i))
}

func (b *Bolt) SetInt(bucket, key []byte, v int) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bkt := tx.Bucket(bucket)
		if err := bkt.Put(key, itob(v)); err != nil {
			return err
		}
		return nil
	})
}

func (b *Bolt) SetString(bucket, key []byte, v string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bkt := tx.Bucket(bucket)
		if err := bkt.Put(key, []byte(v)); err != nil {
			return err
		}
		return nil
	})
}

func itob(i int) []byte {
	b := make([]byte, 32)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
