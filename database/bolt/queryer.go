package bolt

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/index"
)

func (b *Bolt) QueryOne(q index.Query) (index.Result, error) {
	q.Limit = 1
	results, err := b.Query(q)
	if err != nil {
		return index.Result{}, err
	}

	var h index.Hash
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

func (b *Bolt) Query(q index.Query) (index.Results, error) {
	indexVersion, err := b.GetNodeId()
	if err != nil {
		return index.Results{}, errors.Wrap(err, "failed to get index version (nodeid)")
	}

	if q.IndexVersion != "" && indexVersion != q.IndexVersion {
		return index.Results{}, index.ErrIndexVersionsDoNotMatch
	}

	if q.Limit == 0 {
		return index.Results{}, index.ErrNoQueryResults
	}

	if q.FromEntry != 0 {
		hashes := make([]index.Hash, q.Limit)

		var i int
		for ; i < q.Limit; i++ {
			h, err := b.GetIndexEntry(q.FromEntry + i)
			if err != nil && err != database.ErrNoRecord {
				return index.Results{}, errors.Wrap(err, "failed to get entry from db")
			}

			if err == database.ErrNoRecord {
				break
			}

			hashes[i] = index.Hash{
				Entry: q.FromEntry + i,
				Hash:  h,
			}
		}

		return index.Results{
			IndexVersion: indexVersion,
			// trim the slice to the last actual index we got from the db.
			// Ie, if the limit was 50, but only 10 records existed, the slice will be
			// 50 elements big. So indexEntries[:i] will equal indexEntries[:9]
			Hashes: hashes,
		}, nil
	}

	return index.Results{}, index.ErrNoQueryResults
}
