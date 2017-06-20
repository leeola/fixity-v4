package local

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
)

type Blockchain struct {
	lock  *sync.RWMutex
	db    *bolt.DB
	store fixity.Store
	log   log15.Logger
}

func NewBlockchain(log log15.Logger, db *bolt.DB, s fixity.Store) *Blockchain {
	return &Blockchain{
		lock:  &sync.RWMutex{},
		db:    db,
		store: s,
		log:   log,
	}
}

func (l *Blockchain) setHead(h string) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(blockMetaBucketKey)
		if err != nil {
			return err
		}

		return bkt.Put(lastBlockKey, []byte(h))
	})
}

func (l *Blockchain) getHead() (fixity.Block, error) {
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
	if err != nil {
		return fixity.Block{}, err
	}

	if h == "" {
		return fixity.Block{}, fixity.ErrEmptyBlockchain
	}

	var b fixity.Block
	if err := ReadAndUnmarshal(l.store, h, &b); err != nil {
		return fixity.Block{}, err
	}

	b.Hash = h

	return b, nil
}

func (b *Blockchain) AppendContent(c fixity.Content) (fixity.Block, error) {
	if c.Hash == "" {
		return fixity.Block{}, errors.New("content missing Hash value")
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	previousBlock, err := b.getHead()
	if err != nil && err != fixity.ErrEmptyBlockchain {
		return fixity.Block{}, err
	}

	// if the previous hash is the same as the current hash, don't write a new block.
	//
	// There has been no change in the content, so why make a new block.
	if previousBlock.ContentHash == c.Hash {
		b.log.Debug("ignoring identical block",
			"block", previousBlock.Hash,
			"contentHash", c.Hash)
		return previousBlock, nil
	}

	block := fixity.Block{
		// zero value is okay for both of these.
		Block:             previousBlock.Block + 1,
		PreviousBlockHash: previousBlock.Hash,
		ContentHash:       c.Hash,
	}

	bHash, err := MarshalAndWrite(b.store, block)
	if err != nil {
		return fixity.Block{}, err
	}
	block.Hash = bHash
	block.Store = b.store

	return block, nil
}

func (l *Blockchain) Head() (fixity.Block, error) {
	b, err := l.getHead()
	if err != nil {
		return fixity.Block{}, err
	}

	b.Store = l.store
	return b, err
}
