package boltindex

import (
	"github.com/boltdb/bolt"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (bi *BoltIndex) QueryOne(q index.Query) (index.Result, error) {
	q.Limit = 1
	results, err := bi.Query(q)
	if err != nil {
		return index.Result{}, err
	}

	var h string
	// technically Query() should have returned ErrNoQueryResults and been
	// returned above, so there should always be at least one hash. Nevertheless,
	// prevent a slice bounds panic.
	if len(results.Hashes) > 0 {
		h = results.Hashes[0]
	}

	return index.Result{
		IndexVersion: results.IndexVersion,
		Hash:         h,
	}, nil
}

func (bi *BoltIndex) Query(q index.Query) (index.Results, error) {
	if q.IndexVersion != "" && bi.version != q.IndexVersion {
		return index.Results{}, index.ErrIndexVersionsDoNotMatch
	}

	if q.Limit == 0 {
		return index.Results{}, index.ErrNoQueryResults
	}

	if q.FromEntry != 0 {
		indexEntries := make([]string, q.Limit)

		var i int
		for ; i < q.Limit; i++ {
			h, err := bi.GetEntry(q.FromEntry + i)
			if err != nil {
				return index.Results{}, errors.Wrap(err, "failed to get entry from db")
			}
			indexEntries[i] = h
			// db returns zero values for no match.
			if h == "" {
				break
			}
		}

		return index.Results{
			IndexVersion: bi.version,
			// trim the slice to the last actual index we got from the db.
			// Ie, if the limit was 50, but only 10 records existed, the slice will be
			// 50 elements big. So indexEntries[:i] will equal indexEntries[:9]
			Hashes: indexEntries[:i],
		}, nil
	}

	return index.Results{}, index.ErrNoQueryResults
}

func (bi *BoltIndex) Version() string {
	return bi.version
}

func (bi *BoltIndex) AddEntry(h string) error {
	bi.entryCount += 1
	return bi.db.Update(func(tx *bolt.Tx) error {
		entryKey := itob(bi.entryCount)

		metaBucket := tx.Bucket(metaBucketName)
		if err := metaBucket.Put(metaEntryCountKey, entryKey); err != nil {
			return err
		}

		entryBucket := tx.Bucket(entryBucketName)
		if err := entryBucket.Put(entryKey, []byte(h)); err != nil {
			return err
		}

		return nil
	})
}
