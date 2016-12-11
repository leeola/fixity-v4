package blev

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (b *Bleve) QueryOne(q index.Query) (index.Result, error) {
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

func (b *Bleve) Query(q index.Query) (index.Results, error) {
	queries := []query.Query{}

	if q.FromEntry != 0 {
		min, max := float64(q.FromEntry), float64(q.FromEntry+q.Limit)
		nQ := bleve.NewNumericRangeQuery(&min, &max)
		nQ.SetField("index")
		queries = append(queries, nQ)
	}

	if q.Metadata != nil {
		for k, v := range q.Metadata {
			// TODO(leeola): implement meaning method for taking non-string
			s, ok := v.(string)
			if !ok {
				b.log.Warn("unhandled non-string metadata query", "key", k, "value", v)
				continue
			}

			tQ := bleve.NewTermQuery(s)
			tQ.SetField(k)
			queries = append(queries, tQ)
		}
	}

	if len(queries) == 0 {
		return index.Results{
			IndexVersion: b.indexVersion,
		}, nil
	}

	conjQuery := bleve.NewConjunctionQuery(queries...)
	search := bleve.NewSearchRequest(conjQuery)
	search.Fields = []string{"index"}
	searchResults, err := b.index.Search(search)
	if err != nil {
		return index.Results{}, errors.Stack(err)
	}

	hashes := make([]index.Hash, len(searchResults.Hits))

	for i, documentMatch := range searchResults.Hits {
		entryInterface, ok := documentMatch.Fields["index"]
		if !ok {
			return index.Results{}, errors.New("entryIndex value not found")
		}

		entry, ok := entryInterface.(float64)
		if !ok {
			return index.Results{}, errors.New("entryIndex value not int")
		}

		hashes[i] = index.Hash{
			Entry: int(entry),
			Hash:  documentMatch.ID,
		}
	}

	return index.Results{
		IndexVersion: b.indexVersion,
		Hashes:       hashes,
	}, nil
}
