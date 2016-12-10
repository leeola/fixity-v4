package bolt

import (
	boltdb "github.com/boltdb/bolt"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

// TODO(leeola): pull the keys first and append the hash (if not already in).
// This allows multiple keys to be in the query.
func (b *Bolt) Metadata(h string, i index.Indexable) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		metadataBucket := tx.Bucket(metadataBucketName)

		for k, v := range i.ToMetadata() {
			// map the key+value to the db.
			key := k + "::" + v
			err := metadataBucket.Put([]byte(key), []byte(h))
			if err != nil {
				return errors.Wrapf(err, "failed to insert metadata key %q into bolt", key)
			}
		}

		return nil
	})
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
