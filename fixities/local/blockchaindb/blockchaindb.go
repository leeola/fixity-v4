package blockchaindb

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
	Append(fixity.Block) error
	Head() (string, error)
	HeadContentBlock(id string) (string, error)
}

// MutableBlock stores mutable values optimizing the immutable blockchain.
type MutableBlock struct {
	Block                int    `json:"block"`
	Hash                 string `json:"hash"`
	ContentId            string `json:"contentId"`
	PreviousContentBlock int    `json:"previousContentBlock"`
}

type BlockchainDb struct {
	db *bolt.DB
}

func NewBlockchainDb(rootPath string) (*BlockchainDb, error) {
	dbPath := filepath.Join(rootPath, "local", "blockchain.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &BlockchainDb{
		db: db,
	}, nil
}
