package local

import (
	"github.com/boltdb/bolt"
	"github.com/leeola/fixity"
)

type Blockchain struct {
	db    *bolt.DB
	store fixity.Store
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

func (l *Blockchain) getHead() (string, fixity.Block, error) {
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
		return "", fixity.Block{}, err
	}

	if h == "" {
		return "", fixity.Block{}, fixity.ErrEmptyBlockchain
	}

	var b fixity.Block
	if err := ReadAndUnmarshal(l.store, h, &b); err != nil {
		return "", fixity.Block{}, err
	}

	return h, b, nil
}

func (l *Blockchain) Head() (fixity.Block, error) {
	h, b, err := l.getHead()
	if err != nil {
		return fixity.Block{}, err
	}

	b.BlockHash = h
	b.Store = l.store
	return b, err

}
