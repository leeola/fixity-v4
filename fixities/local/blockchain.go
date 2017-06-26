package local

import (
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/fixity"
)

type Blockchain struct {
	lock  *sync.RWMutex
	db    Db
	store fixity.Store
	log   log15.Logger
}

func NewBlockchain(log log15.Logger, db Db, s fixity.Store) *Blockchain {
	return &Blockchain{
		lock:  &sync.RWMutex{},
		db:    db,
		store: s,
		log:   log,
	}
}

func (l *Blockchain) getHead() (fixity.Block, error) {
	h, err := l.db.GetBlockHead()
	if err != nil {
		return fixity.Block{}, err
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

	if err := bc.db.SetBlockHead(bHash); err != nil {
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
