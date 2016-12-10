package bolt

import (
	boltdb "github.com/boltdb/bolt"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (b *Bolt) Metadata(h string, m index.Metadata) error {
	return errors.New("not implemented")
}

func (b *Bolt) Entry(h string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		// get the meta bucket so we can get the total entry count
		metaBkt := tx.Bucket(indexMetaBucketName)

		var entryCount int
		entryCountB := metaBkt.Get(indexMetaEntryCountKey)
		if entryCountB != nil {
			entryCount = btoi(entryCountB)
		}

		// Now that we have the entry count, increment it by one. This will
		// be the key for the hash insertion.
		entryCount += 1
		entryCountB = itob(entryCount)

		// Get the bucket of map[indexes]entries.
		entryBucket := tx.Bucket(indexEntryBucketName)

		// put the map[entryCount]hash into the db.
		err := entryBucket.Put(entryCountB, []byte(h))
		if err != nil {
			return errors.Wrap(err, "failed to insert entry into bolt")
		}

		// Now finally, if that succeeded add the new entry count to the db.
		if err := metaBkt.Put(indexMetaEntryCountKey, entryCountB); err != nil {
			return errors.Wrap(err, "failed to increase entry count metadata")
		}

		return nil
	})
}
