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

func (bc *Blockchain) writeBlock(cb *fixity.ContentBlock, db *fixity.DeleteBlock) (fixity.Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	pb, err := bc.getHead()
	if err != nil && err != fixity.ErrEmptyBlockchain {
		return fixity.Block{}, err
	}

	b := fixity.Block{
		// zero value is okay for both of these.
		Block:             pb.Block + 1,
		PreviousBlockHash: pb.Hash,
		ContentBlock:      cb,
		DeleteBlock:       db,
	}

	bHash, err := MarshalAndWrite(bc.store, b)
	if err != nil {
		return fixity.Block{}, err
	}
	b.Hash = bHash
	b.Store = bc.store

	if err := bc.setHead(bHash); err != nil {
		return fixity.Block{}, err
	}

	return b, nil
}

func (b *Blockchain) AppendContent(c fixity.Content) (fixity.Block, error) {
	if c.Hash == "" {
		return fixity.Block{}, errors.New("Content missing Hash value")
	}

	contentBlock := &fixity.ContentBlock{
		Hash: c.Hash,
	}

	return b.writeBlock(contentBlock, nil)
}

func (bc *Blockchain) DeleteContent(cs ...fixity.Content) (fixity.Block, error) {
	if len(cs) == 0 {
		return fixity.Block{}, nil
	}

	var blocksToBeDeleted []string
	for b, err := bc.Head(); err != fixity.ErrNoMore; b, err = b.Previous() {
		if err != nil {
			return fixity.Block{}, err
		}

		// if this block is not a content block, skip it.
		if b.ContentBlock == nil {
			continue
		}

		for _, c := range cs {
			if b.ContentBlock.Hash == c.Hash {
				blocksToBeDeleted = append(blocksToBeDeleted, b.Hash)
			}
		}
	}

	deleteBlock := &fixity.DeleteBlock{
		Hashes: blocksToBeDeleted,
	}

	return bc.writeBlock(nil, deleteBlock)
}

func (l *Blockchain) Head() (fixity.Block, error) {
	b, err := l.getHead()
	if err != nil {
		return fixity.Block{}, err
	}

	b.Store = l.store
	return b, err
}
