package local

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/leeola/fixity"
)

type Db interface {
	io.Closer

	GetIdHash(id string) (string, error)
	SetIdHash(id, hash string) error
	GetBlockHead() (hash string, err error)
	SetBlockHead(hash string) error
}

type boltDb struct {
	Db *bolt.DB
}

func newBoltDb(rootPath string) (*boltDb, error) {
	dbPath := filepath.Join(rootPath, "local", "local.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &boltDb{Db: db}, nil
}

func (b *boltDb) Close() error {
	return b.Db.Close()
}

func (b *boltDb) GetIdHash(id string) (string, error) {
	var h string
	err := b.Db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(idsBucketKey)
		// if bucket does not exist, this will be nil
		if bkt == nil {
			return nil
		}

		hB := bkt.Get([]byte(id))
		if hB != nil {
			h = string(hB)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if h == "" {
		return "", fixity.ErrIdNotFound
	}

	return h, nil
}

func (b *boltDb) SetIdHash(id, h string) error {
	return b.Db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(idsBucketKey)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(id), []byte(h))
	})
}

func (b *boltDb) GetBlockHead() (string, error) {
	var h string
	err := b.Db.View(func(tx *bolt.Tx) error {
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
	if err != nil {
		return "", err
	}

	if h == "" {
		return "", fixity.ErrNoMore
	}

	return h, nil
}

func (b *boltDb) SetBlockHead(h string) error {
	return b.Db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(blockMetaBucketKey)
		if err != nil {
			return err
		}

		return bkt.Put(lastBlockKey, []byte(h))
	})
}

type memoryDb struct {
	head string
	m    map[string]string
}

func newMemoryDb() *memoryDb {
	return &memoryDb{m: map[string]string{}}
}

func (m *memoryDb) Close() error {
	return nil
}

func (m *memoryDb) GetIdHash(id string) (string, error) {
	h, ok := m.m[id]
	if !ok {
		return "", fixity.ErrIdNotFound
	}
	return h, nil
}

func (m *memoryDb) SetIdHash(id, hash string) error {
	m.m[id] = hash
	return nil
}

func (m *memoryDb) GetBlockHead() (string, error) {
	if m.head == "" {
		return "", fixity.ErrNoMore
	}
	return m.head, nil
}

func (m *memoryDb) SetBlockHead(hash string) error {
	m.head = hash
	return nil
}
